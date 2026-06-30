import { ArrowLeft } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { Separator } from "@/components/ui/separator";
import type {
	Source,
	StepDiff,
	StepResult,
	Workflow,
	WorkflowResult,
} from "@/types";
import { StepCard } from "./StepCard";

interface DetailProps {
	workflow: Workflow;
	sources: Source[];
	diffs: StepDiff[];
	result: WorkflowResult | null;
	mode: "view" | "applying";
	busy: boolean;
	error: string;
	onEdit: () => void;
	onRun: () => void;
	onPreview: () => void;
	onDelete: () => void;
	onPreviewAgain: () => void;
	onBack: () => void;
}

export function WorkflowDetail({
	workflow,
	sources,
	diffs,
	result,
	mode,
	busy,
	error,
	onEdit,
	onRun,
	onPreview,
	onDelete,
	onPreviewAgain,
	onBack,
}: DetailProps) {
	const [runConfirm, setRunConfirm] = useState(false);
	const [deleteConfirm, setDeleteConfirm] = useState(false);

	const diffById: Record<string, StepDiff> = Object.fromEntries(
		diffs.map((d) => [d.step_id, d]),
	);
	const resultById: Record<string, StepResult> = Object.fromEntries(
		(result?.steps ?? []).map((r) => [r.step_id, r]),
	);
	return (
		<div className="flex flex-col h-full">
			{/* ── toolbar ── */}
			<div className="flex h-14 items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10">
				<Button type="button" variant="ghost" size="sm" onClick={onBack}>
					<ArrowLeft className="h-4 w-4 mr-1.5" /> Workflows
				</Button>
				<Separator orientation="vertical" className="h-5 mx-1" />
				<span className="text-sm font-semibold">{workflow.name}</span>
				<div className="flex-1" />
				{error && (
					<span className="text-xs text-destructive mr-2 max-w-xs truncate">
						{error}
					</span>
				)}
				{busy && (
					<span className="text-xs text-muted-foreground mr-2">Loading…</span>
				)}
				{/* Edit */}
				<Button
					type="button"
					variant="outline"
					size="sm"
					onClick={onEdit}
					disabled={busy}
				>
					Edit
				</Button>
				{/* Preview */}
				<Button
					type="button"
					variant="outline"
					size="sm"
					onClick={onPreview}
					disabled={busy}
				>
					Preview
				</Button>
				{/* Run / Preview Again */}
				{mode === "applying" && result ? (
					<Button
						type="button"
						variant="outline"
						size="sm"
						onClick={onPreviewAgain}
						disabled={busy}
					>
						Preview Again
					</Button>
				) : (
					<Button
						type="button"
						size="sm"
						onClick={() => setRunConfirm(true)}
						disabled={busy}
					>
						▶ Run
					</Button>
				)}
				{/* Delete */}
				<Button
					type="button"
					variant="outline"
					size="sm"
					onClick={() => setDeleteConfirm(true)}
					disabled={busy}
					className="text-destructive border-destructive/40 hover:bg-destructive/10"
				>
					Delete
				</Button>
			</div>

			{/* ── steps ── */}
			<div className="flex-1 overflow-auto p-6">
				<div className="flex flex-col gap-4">
					{workflow.steps.length === 0 && (
						<p className="text-sm text-muted-foreground italic">
							No steps. Press Edit to add some.
						</p>
					)}
					{workflow.steps.map((step, i) => (
						<StepCard
							key={step.id || `step-${i}`}
							step={step}
							index={i}
							sources={sources}
							diff={diffById[step.id]}
							result={resultById[step.id]}
							showResult={mode === "applying"}
						/>
					))}
				</div>
			</div>

			{/* ── Run confirmation ── */}
			<ConfirmDialog
				open={runConfirm}
				title="Run this workflow?"
				description="Changes will be applied to your library. This cannot be undone."
				confirmLabel="Run"
				onConfirm={() => {
					setRunConfirm(false);
					onRun();
				}}
				onCancel={() => setRunConfirm(false)}
			/>

			{/* ── Delete confirmation ── */}
			<ConfirmDialog
				open={deleteConfirm}
				title={`Delete "${workflow.name}"?`}
				description="This workflow will be permanently removed."
				confirmLabel="Delete"
				destructive
				onConfirm={() => {
					setDeleteConfirm(false);
					onDelete();
				}}
				onCancel={() => setDeleteConfirm(false)}
			/>
		</div>
	);
}
