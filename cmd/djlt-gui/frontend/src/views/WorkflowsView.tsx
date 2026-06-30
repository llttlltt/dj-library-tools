import { ChevronRight, Eye, PlayCircle, Plus, Trash2 } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import type { QueryTesterOpts } from "@/App";
import { Button } from "@/components/ui/button";
import { Card, CardHeader } from "@/components/ui/card";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { WorkflowDetail } from "@/components/workflow/WorkflowDetail";
import { WorkflowEditor } from "@/components/workflow/WorkflowEditor";
import type { Source, StepDiff, Workflow, WorkflowResult } from "@/types";
import {
	DeleteWorkflow,
	GetWorkflow,
	GetWorkflowDiff,
	ListSources,
	ListWorkflows,
	RunWorkflow,
	SaveWorkflow,
} from "../../wailsjs/go/gui/App";

type Mode = "list" | "view" | "edit" | "applying";

const asWorkflows = (x: unknown) => (x ?? []) as Workflow[];
const asSources = (x: unknown) => (x ?? []) as Source[];
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
	const [wfList, setWfList] = useState<Workflow[]>([]);
	const [selected, setSelected] = useState<Workflow | null>(null);
	const [sources, setSources] = useState<Source[]>([]);
	const [diffs, setDiffs] = useState<StepDiff[]>([]);
	const [result, setResult] = useState<WorkflowResult | null>(null);
	const [error, setError] = useState("");
	const [busy, setBusy] = useState(false);
	const [runTarget, setRunTarget] = useState<Workflow | null>(null);
	const [deleteTarget, setDeleteTarget] = useState<Workflow | null>(null);

	const load = useCallback(async () => {
		try {
			const [wfs, srcs] = await Promise.all([ListWorkflows(), ListSources()]);
			setWfList(asWorkflows(wfs));
			setSources(asSources(srcs));
		} catch (e) {
			setError(String(e));
		}
	}, []);

	useEffect(() => {
		load();
	}, [load]);

	async function openWorkflow(w: Workflow) {
		setError("");
		setDiffs([]);
		setResult(null);
		try {
			const full = asWorkflow(await GetWorkflow(w.id));
			setSelected(JSON.parse(JSON.stringify(full)));
			setMode("view");
		} catch (e) {
			setError(String(e));
		}
	}

	async function openPreview(w: Workflow) {
		setError("");
		setDiffs([]);
		setResult(null);
		try {
			const full = asWorkflow(await GetWorkflow(w.id));
			setSelected(JSON.parse(JSON.stringify(full)));
			setMode("view");
			setBusy(true);
			setDiffs(asDiffs(await GetWorkflowDiff(full.id)));
		} catch (e) {
			setError(String(e));
		}
		setBusy(false);
	}

	async function confirmRun(w: Workflow) {
		setRunTarget(null);
		setError("");
		setDiffs([]);
		setResult(null);
		try {
			const full = asWorkflow(await GetWorkflow(w.id));
			setSelected(JSON.parse(JSON.stringify(full)));
			setMode("applying");
			setBusy(true);
			setResult(asResult(await RunWorkflow(full.id)));
		} catch (e) {
			setError(String(e));
		}
		setBusy(false);
	}

	async function fetchDiff(id: string) {
		setBusy(true);
		setError("");
		try {
			setDiffs(asDiffs(await GetWorkflowDiff(id)));
		} catch (e) {
			setError(String(e));
		}
		setBusy(false);
	}

	async function handleNew() {
		setBusy(true);
		setError("");
		try {
			const wf = asWorkflow(
				await SaveWorkflow({
					id: "",
					name: "New Workflow",
					steps: [],
				} as never),
			);
			await load();
			setSelected(wf);
			setDiffs([]);
			setResult(null);
			setMode("edit");
		} catch (e) {
			setError(String(e));
		}
		setBusy(false);
	}

	async function handleDelete(id: string) {
		setDeleteTarget(null);
		try {
			await DeleteWorkflow(id);
			await load();
			backToList();
		} catch (e) {
			setError(String(e));
		}
	}

	async function handleSave(wf: Workflow) {
		setBusy(true);
		setError("");
		try {
			const saved = asWorkflow(await SaveWorkflow(wf as never));
			setSelected(saved);
			await load();
			setDiffs([]);
			setResult(null);
			setMode("view");
		} catch (e) {
			setError(String(e));
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
			<div className="p-6">
				<div className="flex items-center justify-between mb-6">
					<div>
						<h1 className="text-lg font-semibold">Workflows</h1>
						<p className="text-sm text-muted-foreground mt-0.5">
							Automate your library maintenance
						</p>
					</div>
					<Button size="sm" onClick={handleNew} disabled={busy}>
						<Plus className="h-4 w-4 mr-1.5" /> New Workflow
					</Button>
				</div>
				{error && <p className="text-sm text-destructive mb-4">{error}</p>}
				{wfList.length === 0 ? (
					<p className="text-sm text-muted-foreground italic">
						No workflows yet.
					</p>
				) : (
					<div className="flex flex-col gap-2">
						{wfList.map((w) => (
							<Card
								key={w.id}
								className="cursor-pointer hover:border-border/80 transition-colors"
								onClick={() => openWorkflow(w)}
							>
								<CardHeader className="flex-row items-center justify-between py-3 px-4 gap-0">
									<div className="flex items-center gap-3 min-w-0 flex-1 mr-2">
										<span className="text-sm font-medium truncate">
											{w.name}
										</span>
										<span className="text-xs text-muted-foreground shrink-0">
											{w.steps?.length ?? 0} step
											{w.steps?.length !== 1 ? "s" : ""}
										</span>
									</div>
									<div className="flex items-center gap-0.5">
										<Button
											type="button"
											variant="ghost"
											size="icon"
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
											title="Delete"
											onClick={(e) => {
												e.stopPropagation();
												setDeleteTarget(w);
											}}
										>
											<Trash2 className="h-3.5 w-3.5 text-muted-foreground" />
										</Button>
										<ChevronRight className="h-4 w-4 text-muted-foreground ml-1" />
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
		);

	// ── EDIT ─────────────────────────────────────────────────────────────────
	if (mode === "edit" && selected)
		return (
			<WorkflowEditor
				workflow={selected}
				sources={sources}
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
