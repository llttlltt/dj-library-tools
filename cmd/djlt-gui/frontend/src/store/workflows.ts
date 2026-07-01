import { Atom } from "@effect-atom/atom";
import { Effect } from "effect";
import type { Workflow } from "@/types";
import {
	DeleteWorkflow,
	ListWorkflows,
	SaveWorkflow,
} from "../../wailsjs/go/gui/App";

const asWorkflows = (x: unknown) => (x ?? []) as Workflow[];

// --- State ---

export const workflowsAtom = Atom.make<Workflow[]>([]);
export const workflowsLoadingAtom = Atom.make(false);
export const workflowsErrorAtom = Atom.make<string | null>(null);

// --- Operations ---

export const loadWorkflows = Effect.gen(function* () {
	yield* Atom.set(workflowsLoadingAtom, true);
	yield* Atom.set(workflowsErrorAtom, null);
	try {
		const data = yield* Effect.promise(() => ListWorkflows());
		yield* Atom.set(workflowsAtom, asWorkflows(data));
	} catch (e) {
		yield* Atom.set(workflowsErrorAtom, String(e));
	} finally {
		yield* Atom.set(workflowsLoadingAtom, false);
	}
});

export const saveWorkflow = (wf: Workflow) =>
	Effect.gen(function* () {
		yield* Atom.set(workflowsLoadingAtom, true);
		try {
			yield* Effect.promise(() => SaveWorkflow(wf as never));
			yield* loadWorkflows;
		} catch (e) {
			yield* Atom.set(workflowsErrorAtom, String(e));
		} finally {
			yield* Atom.set(workflowsLoadingAtom, false);
		}
	});

export const removeWorkflow = (id: string) =>
	Effect.gen(function* () {
		yield* Atom.set(workflowsLoadingAtom, true);
		try {
			yield* Effect.promise(() => DeleteWorkflow(id));
			yield* loadWorkflows;
		} catch (e) {
			yield* Atom.set(workflowsErrorAtom, String(e));
		} finally {
			yield* Atom.set(workflowsLoadingAtom, false);
		}
	});
