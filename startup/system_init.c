// System clock initialization - bare metal
typedef unsigned int uint32_t;

// System clock initialization
static void system_init(void) {
    // Configure PLL, clocks, etc.
    // MCU-specific initialization
}

// Called before main()
void __attribute__((constructor)) early_init(void) {
    system_init();
}
