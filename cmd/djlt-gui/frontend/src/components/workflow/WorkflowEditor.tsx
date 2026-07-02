import { Loader2, Plus } from "lucide-react";
import { useState } from "react";
import type { QueryTesterOpts } from "@/App";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { StepCard } from "@/components/workflow/StepCard";
import type { Connection, ProviderInfo, Step, Workflow } from "@/types";

interface EditorProps {
	workflow: Workflow;
	connections: Connection[];
	providers: ProviderInfo[];
	busy: boolean;
	error: string;
	onSave: (wf: Workflow) => void;
	onCancel: () => void;
	onOpenQueryTester?: (opts?: QueryTesterOpts) => void;
}

function blankStep(connectionId: string, targetId: string): Step {
	return {
		id: "",
		kind: "sync",
		source: { connection_id: connectionId, resource: "tracks", query: "" },
		targets: [{ connection_id: targetId, resource: "playlists", query: "" }],
		after: [],
		options: {},
	};
}

export function WorkflowEditor({
	workflow,
	connections,
	providers,
	busy,
	error,
	onSave,
	onCancel,
	onOpenQueryTester,
}: EditorProps) {
	const [wf, setWf] = useState<Workflow>(() =>
		JSON.parse(JSON.stringify(workflow)),
	);
	const [deleteStepIdx, setDeleteStepIdx] = useState<number | null>(null);
	const [showCancelConfirm, setShowCancelConfirm] = useState(false);

	const isDirty = JSON.stringify(wf) !== JSON.stringify(workflow);

	const handleCancel = () => {
		if (isDirty) {
			setShowCancelConfirm(true);
		} else {
			onCancel();
		}
	};

	const sortedConnections = [...connections].sort((a, b) =>
		a.name.localeCompare(b.name),
	);
	const firstConnectionId = sortedConnections[0]?.id ?? "";

	// For targets, pick the first connection that has a provider with CanWrite capability.
	const firstTargetConnectionId =
		sortedConnections.find((c) => {
			const p = providers.find((prov) => prov.name === c.provider);
			return p?.capabilities.CanWrite ?? true;
		})?.id ?? firstConnectionId;

	const mutSteps = (fn: (steps: Step[]) => Step[]) =>
		setWf((w) => ({ ...w, steps: fn([...w.steps]) }));

	const updStep = (i: number, patch: Partial<Step>) =>
		mutSteps((ss) => {
			ss[i] = { ...ss[i], ...patch };
			return ss;
		});

	return (
		<div className="flex flex-col h-full overflow-hidden">
			{/* Sticky Top Header Nav */}
			<div className="h-14 flex items-center gap-3 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10 shrink-0">
				<input
					className="bg-transparent border-none text-sm font-semibold focus:outline-none w-full placeholder:text-muted-foreground/60"
					value={wf.name}
					onChange={(e) => setWf((w) => ({ ...w, name: e.target.value }))}
					placeholder="Workflow name"
				/>
				<div className="flex-1" />
				{error && (
					<span className="text-xs text-destructive mr-2 max-w-xs truncate font-mono">
						{error}
					</span>
				)}
				<Button
					type="button"
					size="sm"
					onClick={() => onSave(wf)}
					disabled={busy}
					className="min-w-[70px]"
				>
					{busy ? <Loader2 className="h-4 w-4 animate-spin mr-1.5" /> : null}
					{busy ? "Saving…" : "Save"}
				</Button>
				<Button
					type="button"
					variant="outline"
					size="sm"
					onClick={handleCancel}
					disabled={busy}
				>
					Cancel
				</Button>
			</div>

			{/* Scrollable Main Content Box */}
			<div className="flex-1 overflow-auto p-6 bg-background">
				<div className="space-y-4">
					{wf.steps.length === 0 ? (
						<p className="text-sm text-muted-foreground italic py-4 pl-1">
							No steps yet — click "+ Add Step" below to get started.
						</p>
					) : (
						<div className="space-y-4">
							{wf.steps.map((step, si) => (
								<StepCard
									key={`editor-step-${step.id || si}`}
									mode="edit"
									step={step}
									index={si}
									connections={connections}
									providers={providers}
									onChange={(patch) => updStep(si, patch)}
									onDelete={() => setDeleteStepIdx(si)}
									onOpenQueryTester={onOpenQueryTester}
								/>
							))}
						</div>
					)}
					<button
						type="button"
						onClick={() =>
							mutSteps((ss) => [
								...ss,
								blankStep(firstConnectionId, firstTargetConnectionId),
							])
						}
						className="w-full flex items-center justify-center gap-1.5 rounded-xl border border-border py-3.5 text-sm font-medium text-muted-foreground hover:border-blue-700/60 hover:text-blue-400 bg-secondary/5 hover:bg-blue-500/[0.02] transition-all duration-200"
					>
						<Plus className="h-4 w-4" /> Add Step
					</button>
				</div>
			</div>

			<ConfirmDialog
				open={showCancelConfirm}
				title="Discard unsaved changes?"
				description="You have unsaved changes in this workflow. If you cancel, these changes will be lost."
				confirmLabel="Discard"
				destructive
				onConfirm={() => {
					setShowCancelConfirm(false);
					onCancel();
				}}
				onCancel={() => setShowCancelConfirm(false)}
			/>

			<ConfirmDialog
				open={deleteStepIdx !== null}
				title="Remove this step?"
				description="This step and its configuration will be removed from the workflow."
				confirmLabel="Remove"
				destructive
				onConfirm={() => {
					if (deleteStepIdx !== null) {
						mutSteps((ss) => ss.filter((_, j) => j !== deleteStepIdx));
					}
					setDeleteStepIdx(null);
				}}
				onCancel={() => setDeleteStepIdx(null)}
			/>
		</div>
	);
}
