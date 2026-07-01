import { CheckCircle, Eye, Loader2, PlayCircle, Trash2 } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
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
}: DetailProps) {
	const [runConfirm, setRunConfirm] = useState(false);
	const [deleteConfirm, setDeleteConfirm] = useState(false);

	const diffsByStepId: Record<string, StepDiff[]> = {};
	for (const d of diffs) {
		if (!diffsByStepId[d.step_id]) diffsByStepId[d.step_id] = [];
		diffsByStepId[d.step_id].push(d);
	}
	const resultById: Record<string, StepResult> = Object.fromEntries(
		(result?.steps ?? []).map((r) => [r.step_id, r]),
	);

	return (
		<div className="flex flex-col h-full overflow-hidden">
			{/* Sticky Top Header Nav */}
			<div className="flex h-14 items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10 shrink-0">
				<span className="text-sm font-semibold">{workflow.name}</span>
				<div className="flex-1" />

				{error && (
					<span className="text-xs text-destructive mr-2 max-w-xs truncate font-mono">
						{error}
					</span>
				)}
				{busy && (
					<div className="flex items-center gap-1.5 text-xs text-muted-foreground mr-2">
						<Loader2 className="h-3.5 w-3.5 animate-spin" />
						Loading…
					</div>
				)}

				<Button
					type="button"
					variant="outline"
					size="sm"
					onClick={() => setDeleteConfirm(true)}
					disabled={busy}
					className="text-[#f43f5e] border-[#f43f5e]/20 hover:border-[#f43f5e]/40 hover:bg-[#f43f5e]/10 transition-colors"
				>
					<Trash2 className="h-4 w-4 mr-1.5" />
					Delete
				</Button>

				<Button
					type="button"
					variant="outline"
					size="sm"
					onClick={onPreview}
					disabled={busy}
				>
					<Eye className="h-4 w-4 mr-1.5" />
					Preview
				</Button>

				{mode === "applying" && result ? (
					<Button
						type="button"
						variant="outline"
						size="sm"
						onClick={onPreviewAgain}
						disabled={busy}
					>
						<CheckCircle className="h-4 w-4 mr-1.5" />
						Preview Again
					</Button>
				) : (
					<Button
						type="button"
						size="sm"
						onClick={() => setRunConfirm(true)}
						disabled={busy}
					>
						<PlayCircle className="h-4 w-4 mr-1.5" />
						Run
					</Button>
				)}

				<Button
					type="button"
					variant="outline"
					size="sm"
					onClick={onEdit}
					disabled={busy}
				>
					Edit
				</Button>
			</div>

			{/* Scrollable Main Content Box */}
			<div className="flex-1 overflow-auto p-6 bg-background">
				<div className="space-y-4">
					{workflow.steps.length === 0 && (
						<p className="text-sm text-muted-foreground italic pl-1 py-2">
							No steps configured. Press "Edit" above to add some.
						</p>
					)}
					{workflow.steps.map((step, i) => (
						<StepCard
							key={step.id || `step-${i}`}
							mode={mode} // Propagates "view" or "applying"
							step={step}
							index={i}
							sources={sources}
							diffs={diffsByStepId[step.id] || []}
							result={resultById[step.id]}
						/>
					))}
				</div>
			</div>

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
