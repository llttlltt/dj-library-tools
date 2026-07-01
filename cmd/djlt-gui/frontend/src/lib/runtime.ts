import { Registry } from "@effect-atom/atom";
import { ManagedRuntime } from "effect";

// Create a runtime that includes the AtomRegistry service.
export const runtime = ManagedRuntime.make(Registry.layer);

// Extract the singleton Registry instance from the runtime
// so it can be passed to the React Context.
export const registry: Registry.Registry = runtime.runSync(
	Registry.AtomRegistry,
);

// Export helper to run effects that require the AtomRegistry.
export const runPromise = runtime.runPromise;
export const runSync = runtime.runSync;
