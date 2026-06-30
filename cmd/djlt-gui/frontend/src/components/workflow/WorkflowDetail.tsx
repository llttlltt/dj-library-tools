import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import type { StepDiff, StepResult, Workflow, WorkflowResult } from "@/types";
import { StepCard } from "./StepCard";

interface DetailProps {
	workflow: Workflow;
	diffs: StepDiff[];
	result: WorkflowResult | null;
	mode: "view" | "applying";
	busy: boolean;
	error: string;
	onEdit: () => void;
	onApply: () => void;
	onPreview: () => void;
	onDelete: () => void;
	onPreviewAgain: () => void;
	onBack: () => void;
}

export function WorkflowDetail({
	workflow,
	diffs,
	result,
	mode,
	busy,
	error,
	onEdit,
	onApply,
	onPreview,
	onDelete,
	onPreviewAgain,
	onBack,
}: DetailProps) {
	const diffById: Record<string, StepDiff> = Object.fromEntries(
		diffs.map((d) => [d.step_id, d]),
	);
	const resultById: Record<string, StepResult> = Object.fromEntries(
		(result?.steps ?? []).map((r) => [r.step_id, r]),
	);
	const syncSteps = workflow.steps.filter((s) => s.kind === "sync").length;
	const diffLoaded = diffs.length > 0 || syncSteps === 0;

	return (
		<div className="flex flex-col h-full">
			<div className="flex items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10">
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
				<Button
					type="button"
					variant="outline"
					size="sm"
					onClick={onEdit}
					disabled={busy}
				>
					Edit
				</Button>
				<Button
					type="button"
					variant="outline"
					size="sm"
					onClick={onPreview}
					disabled={busy}
				>
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
						Preview Again
					</Button>
				) : (
					<Button
						type="button"
						size="sm"
						onClick={onApply}
						disabled={!diffLoaded || busy}
					>
						▶ Run
					</Button>
				)}
				<Button
					type="button"
					variant="outline"
					size="sm"
					onClick={onDelete}
					disabled={busy}
					className="text-destructive border-destructive/40 hover:bg-destructive/10"
				>
					Delete
				</Button>
			</div>

			<div className="flex-1 overflow-auto p-6">
				<div className="flex flex-col gap-4 max-w-3xl">
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
							diff={diffById[step.id]}
							result={resultById[step.id]}
							showResult={mode === "applying"}
						/>
					))}
				</div>
			</div>
		</div>
	);
}
