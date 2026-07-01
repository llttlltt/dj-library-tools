import { Duration, Effect, Schedule, Schema } from "effect";
import * as WailsApp from "../../wailsjs/go/gui/App";
import {
	WailsCallError,
	WailsDecodeError,
	WailsRuntimeNotReadyError,
} from "./errors";
import {
	ProviderInfoSchema,
	QueryResultSchema,
	SourceSchema,
	StepDiffSchema,
	WorkflowSchema,
} from "./schemas";

/**
 * Structural relaxes a Wails-generated class instance type to a plain object
 * shape, so we can pass POJOs to bindings that nominally expect model classes.
 * This version is recursive to handle nested model arrays (like Step[] in Workflow).
 */
// biome-ignore lint/suspicious/noExplicitAny: generic recursive model mapper
type Structural<T> = T extends (...args: any[]) => any
	? never
	: T extends object
		? {
				// biome-ignore lint/suspicious/noExplicitAny: generic function matcher
				[K in keyof T as T[K] extends (...args: any[]) => any
					? never
					: K]: Structural<T[K]>;
			}
		: T;

/**
 * Confirms the Wails runtime has been injected into the window.
 * Bindings call into `window.go.*`, which is absent during SSR/tests/startup.
 */
const ensureRuntime = Effect.suspend(() =>
	typeof window !== "undefined" && (window as { go?: unknown }).go
		? Effect.void
		: Effect.fail(
				new WailsRuntimeNotReadyError({
					message: "Wails runtime not available (window.go is undefined)",
				}),
			),
);

/**
 * Base logic for all Wails calls.
 */
const wailsCall =
	<Args extends readonly unknown[], A>(fn: (...args: Args) => Promise<A>) =>
	(...args: Args) =>
		ensureRuntime.pipe(
			Effect.andThen(() =>
				Effect.tryPromise({
					try: () => fn(...args),
					catch: (originalError) =>
						new WailsCallError({
							message: String(originalError),
							originalError,
						}),
				}),
			),
			Effect.timeout(Duration.seconds(30)),
		);

/**
 * Resilience policy for READ IPC calls.
 */
const readResilience = Effect.retry({
	schedule: Schedule.exponential(Duration.millis(100)).pipe(
		Schedule.compose(Schedule.recurs(2)),
	),
});

/**
 * AppService defines the Effect-native interface for Wails Go bindings.
 */
export class AppService extends Effect.Service<AppService>()("AppService", {
	accessors: true,
	sync: () => {
		const read =
			// biome-ignore lint/suspicious/noExplicitAny: generic function wrapper
				<Args extends readonly any[], A>(fn: (...args: Args) => Promise<A>) =>
				(...args: Args) =>
					wailsCall(fn)(...args).pipe(readResilience);

		const readDecoded =
			// biome-ignore lint/suspicious/noExplicitAny: generic function wrapper
				<Args extends readonly any[], A, I, R>(
					fn: (...args: Args) => Promise<A>,
					// biome-ignore lint/suspicious/noExplicitAny: allow any context for schema
					schema: Schema.Schema<R, I, any>,
				) =>
				(...args: Args) =>
					Effect.gen(function* () {
						const data = yield* read(fn)(...args);
						return yield* Schema.decodeUnknown(
							schema as Schema.Schema<R, I, never>,
						)(data).pipe(
							Effect.mapError(
								(parseError) =>
									new WailsDecodeError({
										message: "Failed to decode Go payload",
										parseError,
									}),
							),
						);
					});

		const write = wailsCall;

		const writeStructural =
			// biome-ignore lint/suspicious/noExplicitAny: generic function wrapper
				<Args extends readonly any[], A>(fn: (...args: Args) => Promise<A>) =>
				(...args: { [K in keyof Args]: Structural<Args[K]> }) =>
					write(fn)(...(args as unknown as Args));

		return {
			listSources: readDecoded(
				WailsApp.ListSources,
				Schema.Array(SourceSchema),
			),
			listWorkflows: readDecoded(
				WailsApp.ListWorkflows,
				Schema.Array(WorkflowSchema),
			),
			getWorkflow: readDecoded(WailsApp.GetWorkflow, WorkflowSchema),
			getWorkflowDiff: readDecoded(
				WailsApp.GetWorkflowDiff,
				Schema.Array(StepDiffSchema),
			),
			listProviders: readDecoded(
				WailsApp.ListProviders,
				Schema.Array(ProviderInfoSchema),
			),
			previewQuery: readDecoded(WailsApp.PreviewQuery, QueryResultSchema),

			createSource: write(WailsApp.CreateSource),
			deleteSource: write(WailsApp.DeleteSource),
			updateSource: writeStructural(WailsApp.UpdateSource),
			saveWorkflow: writeStructural(WailsApp.SaveWorkflow),
			deleteWorkflow: write(WailsApp.DeleteWorkflow),
			runWorkflow: write(WailsApp.RunWorkflow),
		};
	},
}) {}

export const WailsLive = AppService.Default;
