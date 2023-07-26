#include "_cgo_export.h"
#include "sketchybar.h"
void MainCFunction(char* bootstrapName) {
//	printf("MainCFunction()\n");
//	AGoFunction();

	event_server_begin(handleEvent, bootstrapName);
}

// Define the Go wrapper for the mach_handler function
void handleEvent(env env) {
    // Your Go code to handle the event goes here
    // The 'env' argument can be converted to a Go string using C.GoString
    // Environment variables passed from sketchybar can be accessed as seen below
    char* name = env_get_value_for_key(env, "NAME");
    char* sender = env_get_value_for_key(env, "SENDER");
    char* info = env_get_value_for_key(env, "INFO");
    char* config_dir = env_get_value_for_key(env, "CONFIG_DIR");
    char* button = env_get_value_for_key(env, "BUTTON");
    char* modifier = env_get_value_for_key(env, "MODIFIER");

    GoHandler(name, sender, info, config_dir, button, modifier);
}

