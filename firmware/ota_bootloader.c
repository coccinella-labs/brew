// OTA (Over-The-Air) firmware update system
#include <stdint.h>

#define FLASH_BASE      0x08000000
#define APP_BASE        0x08008000  // App starts at 32KB
#define BOOTLOADER_SIZE 0x8000      // 32KB bootloader
#define APP_SIZE        0x38000     // 224KB app space

// Firmware metadata
typedef struct {
    uint32_t magic;      // 0xDEADBEEF
    uint32_t version;    // Firmware version
    uint32_t size;       // Firmware size
    uint32_t crc32;      // CRC checksum
} firmware_header_t;

// OTA state machine
typedef enum {
    OTA_IDLE,
    OTA_DOWNLOADING,
    OTA_VERIFYING,
    OTA_INSTALLING,
    OTA_COMPLETE
} ota_state_t;

static ota_state_t ota_state = OTA_IDLE;
static uint32_t download_addr = 0x08040000; // Temp storage

void flash_unlock(void) {
    // Unlock flash for writing
    *((volatile uint32_t*)0x40023C04) = 0x45670123;
    *((volatile uint32_t*)0x40023C04) = 0xCDEF89AB;
}

void flash_write_word(uint32_t addr, uint32_t data) {
    *((volatile uint32_t*)addr) = data;
    while (*((volatile uint32_t*)0x40023C0C) & 0x10000); // Wait
}

uint32_t calculate_crc32(uint32_t addr, uint32_t size) {
    // Simple CRC32 implementation
    uint32_t crc = 0xFFFFFFFF;
    uint8_t *data = (uint8_t*)addr;

    for (uint32_t i = 0; i < size; i++) {
        crc ^= data[i];
        for (int j = 0; j < 8; j++) {
            crc = (crc >> 1) ^ (0xEDB88320 & (-(crc & 1)));
        }
    }
    return ~crc;
}

void ota_process_chunk(uint8_t *data, uint32_t offset, uint32_t size) {
    if (ota_state != OTA_DOWNLOADING) return;

    flash_unlock();

    // Write chunk to temporary location
    for (uint32_t i = 0; i < size; i += 4) {
        uint32_t word = *(uint32_t*)(data + i);
        flash_write_word(download_addr + offset + i, word);
    }
}

void ota_verify_and_install(void) {
    firmware_header_t *header = (firmware_header_t*)download_addr;

    // Verify magic and CRC
    if (header->magic != 0xDEADBEEF) {
        ota_state = OTA_IDLE;
        return;
    }

    uint32_t calc_crc = calculate_crc32(download_addr + sizeof(firmware_header_t),
                                       header->size);
    if (calc_crc != header->crc32) {
        ota_state = OTA_IDLE;
        return;
    }

    ota_state = OTA_INSTALLING;

    // Copy verified firmware to app location
    flash_unlock();
    uint32_t *src = (uint32_t*)download_addr;
    uint32_t *dst = (uint32_t*)APP_BASE;

    for (uint32_t i = 0; i < header->size / 4; i++) {
        flash_write_word((uint32_t)(dst + i), src[i]);
    }

    ota_state = OTA_COMPLETE;

    // Reset to run new firmware
    *((volatile uint32_t*)0xE000ED0C) = 0x05FA0004;
}
