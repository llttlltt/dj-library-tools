import { Chunk, Effect, Stream } from "effect";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";

/**
 * Subscribes to a Wails runtime event as an Effect Stream. The subscription is
 * automatically torn down when the stream is interrupted or completes.
 */
export const wailsEvents = <A = unknown>(eventName: string): Stream.Stream<A> =>
	Stream.asyncScoped<A>((emit) =>
		Effect.acquireRelease(
			Effect.sync(() => {
				EventsOn(eventName, (...data: A[]) => {
					// Wails passes variadic payloads; emit each element.
					emit(Effect.succeed(Chunk.fromIterable(data)));
				});
			}),
			() => Effect.sync(() => EventsOff(eventName)),
		),
	);
