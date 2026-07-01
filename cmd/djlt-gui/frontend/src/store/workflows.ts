import { Atom } from "@effect-atom/atom";
import { Effect } from "effect";
import { AppService } from "@/services";
import type { Workflow } from "@/types";

const asWorkflows = (x: unknown) => (x ?? []) as Workflow[];

// --- State ---

export const workflowsAtom = Atom.make<Workflow[]>([]);
export const workflowsLoadingAtom = Atom.make(false);
export const workflowsErrorAtom = Atom.make<string | null>(null);

// --- Operations ---

/**
 * Normalizes any error into a string for the error atom.
 */
const handleError = (e: unknown) =>
	Atom.set(workflowsErrorAtom, e instanceof Error ? e.message : String(e));

export const loadWorkflows = Effect.gen(function* () {
	yield* Atom.set(workflowsLoadingAtom, true);
	yield* Atom.set(workflowsErrorAtom, null);
	const app = yield* AppService;
	return yield* app.listWorkflows().pipe(
		Effect.flatMap((data) => Atom.set(workflowsAtom, asWorkflows(data))),
		Effect.catchAll((e) => handleError(e)),
		Effect.andThen(() => Atom.set(workflowsLoadingAtom, false)),
	);
});

export const saveWorkflow = (wf: Workflow) =>
	Effect.gen(function* () {
		yield* Atom.set(workflowsLoadingAtom, true);
		const app = yield* AppService;
		return yield* app.saveWorkflow(wf).pipe(
			Effect.flatMap(() => loadWorkflows),
			Effect.catchAll((e) => handleError(e)),
			Effect.andThen(() => Atom.set(workflowsLoadingAtom, false)),
		);
	});

export const removeWorkflow = (id: string) =>
	Effect.gen(function* () {
		yield* Atom.set(workflowsLoadingAtom, true);
		const app = yield* AppService;
		return yield* app.deleteWorkflow(id).pipe(
			Effect.flatMap(() => loadWorkflows),
			Effect.catchAll((e) => handleError(e)),
			Effect.andThen(() => Atom.set(workflowsLoadingAtom, false)),
		);
	});
