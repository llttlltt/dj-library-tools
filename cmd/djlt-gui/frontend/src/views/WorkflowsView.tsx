import { useAtom } from "@effect-atom/atom-react";
import { ChevronRight, Eye, PlayCircle, Plus, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";
import type { QueryTesterOpts } from "@/App";
import { Button } from "@/components/ui/button";
import { Card, CardHeader } from "@/components/ui/card";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { WorkflowDetail } from "@/components/workflow/WorkflowDetail";
import { WorkflowEditor } from "@/components/workflow/WorkflowEditor";
import { runtime } from "@/lib/runtime";
import { loadProviders, providersAtom } from "@/store/providers";
import { loadSources, sourcesAtom } from "@/store/sources";
import {
	loadWorkflows,
	removeWorkflow,
	saveWorkflow,
	workflowsAtom,
	workflowsErrorAtom,
} from "@/store/workflows";
import type { StepDiff, Workflow, WorkflowResult } from "@/types";
import {
	GetWorkflow,
	GetWorkflowDiff,
	RunWorkflow,
} from "../../wailsjs/go/gui/App";

type Mode = "list" | "view" | "edit" | "applying";

const asWorkflow = (x: unknown) => x as Workflow;
const asDiffs = (x: unknown) => (x ?? []) as StepDiff[];
const asResult = (x: unknown) => x as WorkflowResult;

interface WorkflowsViewProps {
	onOpenQueryTester: (opts?: QueryTesterOpts) => void;
}

export default function WorkflowsView({
	onOpenQueryTester,
}: WorkflowsViewProps) {
	const [mode, setMode] = useState<Mode>("list");
	const [wfList] = useAtom(workflowsAtom);
	const [errorValue, setError] = useAtom(workflowsErrorAtom);
	const error = errorValue ?? "";
	const [sources] = useAtom(sourcesAtom);
	const [providers] = useAtom(providersAtom);

	const [selected, setSelected] = useState<Workflow | null>(null);
	const [diffs, setDiffs] = useState<StepDiff[]>([]);
	const [result, setResult] = useState<WorkflowResult | null>(null);
	const [busy, setBusy] = useState(false);
	const [runTarget, setRunTarget] = useState<Workflow | null>(null);
	const [deleteTarget, setDeleteTarget] = useState<Workflow | null>(null);

	useEffect(() => {
		runtime.runPromise(loadWorkflows);
		runtime.runPromise(loadSources);
		runtime.runPromise(loadProviders);
	}, []);

	async function openWorkflow(w: Workflow) {
		setDiffs([]);
		setResult(null);
		try {
			const full = asWorkflow(await GetWorkflow(w.id));
			setSelected(JSON.parse(JSON.stringify(full)));
			setMode("view");
		} catch (e) {
			console.error(e);
		}
	}

	async function openPreview(w: Workflow) {
		setDiffs([]);
		setResult(null);
		try {
			const full = asWorkflow(await GetWorkflow(w.id));
			setSelected(JSON.parse(JSON.stringify(full)));
			setMode("view");
			setBusy(true);
			setDiffs(asDiffs(await GetWorkflowDiff(full.id)));
		} catch (e) {
			console.error(e);
		}
		setBusy(false);
	}

	async function confirmRun(w: Workflow) {
		setRunTarget(null);
		setDiffs([]);
		setResult(null);
		try {
			const full = asWorkflow(await GetWorkflow(w.id));
			setSelected(JSON.parse(JSON.stringify(full)));
			setMode("applying");
			setBusy(true);
			setResult(asResult(await RunWorkflow(full.id)));
		} catch (e) {
			console.error(e);
		}
		setBusy(false);
	}

	async function fetchDiff(id: string) {
		setBusy(true);
		try {
			setDiffs(asDiffs(await GetWorkflowDiff(id)));
		} catch (e) {
			console.error(e);
		}
		setBusy(false);
	}

	async function handleNew() {
		setBusy(true);
		try {
			const wf = {
				id: "",
				name: "New Workflow",
				steps: [],
			} as Workflow;

			await runtime.runPromise(saveWorkflow(wf));
			// The store update will trigger a re-render, but we need to find the new ID.
			// Actually, saveWorkflow already calls loadWorkflows.
			// For simplicity in this first pass, we'll just go back to list or stay in new state.
			setMode("list");
		} catch (e) {
			console.error(e);
		}
		setBusy(false);
	}

	async function handleDelete(id: string) {
		setDeleteTarget(null);
		await runtime.runPromise(removeWorkflow(id));
		backToList();
	}

	async function handleSave(wf: Workflow) {
		setBusy(true);
		try {
			await runtime.runPromise(saveWorkflow(wf));
			// Update the selected workflow to reflect the saved state
			setSelected(JSON.parse(JSON.stringify(wf)));
			setMode("view");
			// Clear stale diffs/results
			setDiffs([]);
			setResult(null);
		} catch (e) {
			console.error(e);
		}
		setBusy(false);
	}

	async function handleRun() {
		if (!selected) return;
		setBusy(true);
		setError("");
		setMode("applying");
		try {
			setResult(asResult(await RunWorkflow(selected.id)));
		} catch (e) {
			setError(String(e));
		}
		setBusy(false);
	}

	function backToList() {
		setMode("list");
		setSelected(null);
		setDiffs([]);
		setResult(null);
		setError("");
	}

	// ── LIST ─────────────────────────────────────────────────────────────────
	if (mode === "list")
		return (
			<div className="flex flex-col h-full">
				<div className="h-14 flex items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10">
					<span className="text-sm font-semibold">Workflows</span>
					<div className="flex-1" />
					<Button size="sm" onClick={handleNew} disabled={busy}>
						<Plus className="h-4 w-4 mr-1.5" /> New Workflow
					</Button>
				</div>
				<div className="flex-1 overflow-auto p-6">
					{error && <p className="text-sm text-destructive mb-4">{error}</p>}

					{wfList.length === 0 ? (
						<p className="text-sm text-muted-foreground italic">
							No workflows yet. Click "New Workflow" to get started.
						</p>
					) : (
						<div className="flex flex-col gap-3">
							{[...wfList]
								.sort((a, b) => a.name.localeCompare(b.name))
								.map((w) => (
									<Card
										key={w.id}
										className="cursor-pointer hover:border-border/80 transition-colors"
										onClick={() => openWorkflow(w)}
									>
										<CardHeader className="flex-row items-center justify-between py-3 px-4 gap-0 space-y-0">
											<div className="flex items-center gap-3 min-w-0 flex-1 mr-2">
												<span className="text-sm font-medium truncate">
													{w.name}
												</span>
												<span className="text-xs text-muted-foreground shrink-0 bg-secondary/50 py-0.5 px-2 rounded-full">
													{w.steps?.length ?? 0} step
													{w.steps?.length !== 1 ? "s" : ""}
												</span>
											</div>
											<div className="flex items-center gap-1">
												<Button
													type="button"
													variant="ghost"
													size="icon"
													className="h-8 w-8 hover:bg-secondary"
													title="Run"
													onClick={(e) => {
														e.stopPropagation();
														setRunTarget(w);
													}}
												>
													<PlayCircle className="h-4 w-4 text-emerald-500" />
												</Button>
												<Button
													type="button"
													variant="ghost"
													size="icon"
													className="h-8 w-8 hover:bg-secondary"
													title="Preview"
													onClick={(e) => {
														e.stopPropagation();
														openPreview(w);
													}}
												>
													<Eye className="h-4 w-4 text-blue-400" />
												</Button>
												<Button
													type="button"
													variant="ghost"
													size="icon"
													className="h-8 w-8 hover:bg-secondary"
													title="Delete"
													onClick={(e) => {
														e.stopPropagation();
														setDeleteTarget(w);
													}}
												>
													<Trash2 className="h-4 w-4 text-muted-foreground hover:text-destructive" />
												</Button>
												<div className="h-4 w-px bg-border mx-1" />
												<ChevronRight className="h-4 w-4 text-muted-foreground/60" />
											</div>
										</CardHeader>
									</Card>
								))}
						</div>
					)}

					{/* Run confirm — list level */}
					{runTarget && (
						<ConfirmDialog
							open={true}
							title={`Run "${runTarget.name}"?`}
							description="Changes will be applied to your library. This cannot be undone."
							confirmLabel="Run"
							onConfirm={() => confirmRun(runTarget)}
							onCancel={() => setRunTarget(null)}
						/>
					)}
					{/* Delete confirm — list level */}
					{deleteTarget && (
						<ConfirmDialog
							open={true}
							title={`Delete "${deleteTarget.name}"?`}
							description="This workflow will be permanently removed."
							confirmLabel="Delete"
							destructive
							onConfirm={() => handleDelete(deleteTarget.id)}
							onCancel={() => setDeleteTarget(null)}
						/>
					)}
				</div>
			</div>
		);

	// ── EDIT ─────────────────────────────────────────────────────────────────
	if (mode === "edit" && selected)
		return (
			<WorkflowEditor
				workflow={selected}
				sources={sources}
				providers={providers}
				busy={busy}
				error={error}
				onSave={handleSave}
				onOpenQueryTester={onOpenQueryTester}
				onCancel={() => {
					if (selected.id) {
						setMode("view");
					} else {
						backToList();
					}
				}}
			/>
		);

	// ── VIEW / APPLYING ───────────────────────────────────────────────────────
	if ((mode === "view" || mode === "applying") && selected)
		return (
			<WorkflowDetail
				workflow={selected}
				sources={sources}
				diffs={diffs}
				result={result}
				mode={mode}
				busy={busy}
				error={error}
				onEdit={() => setMode("edit")}
				onRun={handleRun}
				onPreview={() => fetchDiff(selected.id)}
				onDelete={() => handleDelete(selected.id)}
				onPreviewAgain={() => fetchDiff(selected.id)}
				onBack={backToList}
			/>
		);

	return null;
}
