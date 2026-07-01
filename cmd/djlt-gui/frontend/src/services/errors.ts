import type { ParseResult } from "effect";
import { Data } from "effect";

/**
 * Error raised when the Wails runtime is not yet injected into `window`.
 */
export class WailsRuntimeNotReadyError extends Data.TaggedError(
	"WailsRuntimeNotReadyError",
)<{
	readonly message: string;
}> {}

/**
 * Error raised when a Go payload fails Schema decoding at the IPC boundary.
 */
export class WailsDecodeError extends Data.TaggedError("WailsDecodeError")<{
	readonly message: string;
	readonly parseError: ParseResult.ParseError;
}> {}

/**
 * Base class for errors occurring during a Wails binding call.
 */
export class WailsCallError extends Data.TaggedError("WailsCallError")<{
	readonly message: string;
	readonly originalError: unknown;
}> {}
