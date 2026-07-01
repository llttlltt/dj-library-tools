import { Duration, Effect, Schedule, Schema } from "effect";
import * as WailsApp from "../../wailsjs/go/gui/App";
import {
	WailsCallError,
	WailsDecodeError,
	WailsRuntimeNotReadyError,
} from "./errors";
import {
	ConnectionSchema,
	ProviderInfoSchema,
	QueryResultSchema,
	StepDiffSchema,
	UpdateConfigSchema,
	UpdateInfoSchema,
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
			// System
			getVersion: read(WailsApp.GetVersion),
			getUpdateConfig: readDecoded(
				WailsApp.GetUpdateConfig,
				UpdateConfigSchema,
			),
			setUpdateInterval: write(WailsApp.SetUpdateInterval),
			checkForUpdate: readDecoded(WailsApp.CheckForUpdate, UpdateInfoSchema),
			installUpdate: write(WailsApp.InstallUpdate),
			getPermissionStatus: read(WailsApp.GetPermissionStatus),
			fixPermissions: write(WailsApp.FixPermissions),
			openFileDialog: read(WailsApp.OpenFileDialog),

			// Plex Auth
			initPlexAuth: read(WailsApp.InitPlexAuth),
			checkPlexAuth: read(WailsApp.CheckPlexAuth),

			// Connections
			listConnections: readDecoded(
				WailsApp.ListConnections,
				Schema.Array(ConnectionSchema),
			),
			createConnection: write(WailsApp.CreateConnection),
			deleteConnection: write(WailsApp.DeleteConnection),
			updateConnection: writeStructural(WailsApp.UpdateConnection),

			// Workflows
			listWorkflows: readDecoded(
				WailsApp.ListWorkflows,
				Schema.Array(WorkflowSchema),
			),
			getWorkflow: readDecoded(WailsApp.GetWorkflow, WorkflowSchema),
			getWorkflowDiff: readDecoded(
				WailsApp.GetWorkflowDiff,
				Schema.Array(StepDiffSchema),
			),
			saveWorkflow: writeStructural(WailsApp.SaveWorkflow),
			deleteWorkflow: write(WailsApp.DeleteWorkflow),
			runWorkflow: write(WailsApp.RunWorkflow),

			// Providers / Library
			listProviders: readDecoded(
				WailsApp.ListProviders,
				Schema.Array(ProviderInfoSchema),
			),
			previewQuery: readDecoded(WailsApp.PreviewQuery, QueryResultSchema),
		};
	},
}) {}

export const WailsLive = AppService.Default;
