import { Registry } from "@effect-atom/atom";
import { Layer, ManagedRuntime } from "effect";
import { WailsLive } from "@/services/app-service";

/**
 * AppRuntime provides the Effect execution context for the application,
 * including Wails IPC and the Atom registry.
 */
export const AppRuntime = ManagedRuntime.make(
	Layer.merge(Registry.layer, WailsLive),
);

// Compatibility exports for existing store logic
export const runtime = AppRuntime;
export const registry: Registry.Registry = runtime.runSync(
	Registry.AtomRegistry,
);
export const runPromise = runtime.runPromise;
export const runSync = runtime.runSync;
