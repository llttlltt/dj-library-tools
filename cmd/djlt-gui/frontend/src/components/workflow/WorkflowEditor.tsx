import { X } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import type { Endpoint, Source, Step, Workflow } from "@/types";

// ── EpEditRow ──────────────────────────────────────────────────────────────

interface EpEditRowProps {
	ep: Endpoint;
	sources: Source[];
	onChange: (p: Partial<Endpoint>) => void;
}

export function EpEditRow({ ep, sources, onChange }: EpEditRowProps) {
	return (
		<div className="flex gap-2 items-center">
			<Select
				value={ep.source_id}
				onValueChange={(v) => onChange({ source_id: v })}
			>
				<SelectTrigger className="w-36 h-7 text-xs shrink-0">
					<SelectValue placeholder="Source" />
				</SelectTrigger>
				<SelectContent>
					{sources.map((s) => (
						<SelectItem key={s.id} value={s.id}>
							{s.name}
						</SelectItem>
					))}
				</SelectContent>
			</Select>
			<Input
				className="h-7 text-xs w-24 shrink-0"
				value={ep.resource}
				onChange={(e) => onChange({ resource: e.target.value })}
				placeholder="resource"
			/>
			<Input
				className="h-7 text-xs flex-1"
				value={ep.query ?? ""}
				onChange={(e) => onChange({ query: e.target.value })}
				placeholder="query (optional)"
			/>
		</div>
	);
}

// ── WorkflowEditor ─────────────────────────────────────────────────────────

interface EditorProps {
	workflow: Workflow;
	sources: Source[];
	busy: boolean;
	error: string;
	onSave: (wf: Workflow) => void;
	onCancel: () => void;
}

function blankStep(srcId: string): Step {
	return {
		id: "",
		kind: "sync",
		source: { source_id: srcId, resource: "tracks", query: "" },
		targets: [{ source_id: srcId, resource: "playlists", query: "" }],
		after: [],
		options: {},
	};
}

export function WorkflowEditor({
	workflow,
	sources,
	busy,
	error,
	onSave,
	onCancel,
}: EditorProps) {
	const [wf, setWf] = useState<Workflow>(() =>
		JSON.parse(JSON.stringify(workflow)),
	);
	const firstSrcId = sources[0]?.id ?? "";

	const mutSteps = (fn: (steps: Step[]) => Step[]) =>
		setWf((w) => ({ ...w, steps: fn([...w.steps]) }));

	const updStep = (i: number, patch: Partial<Step>) =>
		mutSteps((ss) => {
			ss[i] = { ...ss[i], ...patch };
			return ss;
		});

	const updSource = (si: number, patch: Partial<Endpoint>) =>
		mutSteps((ss) => {
			ss[si] = { ...ss[si], source: { ...ss[si].source, ...patch } };
			return ss;
		});

	const updTarget = (si: number, ti: number, patch: Partial<Endpoint>) =>
		mutSteps((ss) => {
			const tgts = [...ss[si].targets];
			tgts[ti] = { ...tgts[ti], ...patch };
			ss[si] = { ...ss[si], targets: tgts };
			return ss;
		});

	const addTarget = (si: number) =>
		mutSteps((ss) => {
			ss[si] = {
				...ss[si],
				targets: [
					...ss[si].targets,
					{ source_id: firstSrcId, resource: "playlists", query: "" },
				],
			};
			return ss;
		});

	const removeTarget = (si: number, ti: number) =>
		mutSteps((ss) => {
			ss[si] = {
				...ss[si],
				targets: ss[si].targets.filter((_, j) => j !== ti),
			};
			return ss;
		});

	return (
		<div className="flex flex-col h-full">
			<div className="flex items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10">
				<Button
					type="button"
					variant="ghost"
					size="sm"
					onClick={onCancel}
					disabled={busy}
				>
					<X className="h-4 w-4 mr-1" /> Cancel
				</Button>
				<Separator orientation="vertical" className="h-5 mx-1" />
				<input
					className="bg-transparent border-none text-sm font-semibold focus:outline-none w-64"
					value={wf.name}
					onChange={(e) => setWf((w) => ({ ...w, name: e.target.value }))}
					placeholder="Workflow name"
				/>
				<div className="flex-1" />
				{error && (
					<span className="text-xs text-destructive mr-2 max-w-xs truncate">
						{error}
					</span>
				)}
				<Button
					type="button"
					size="sm"
					onClick={() => onSave(wf)}
					disabled={busy}
				>
					{busy ? "Saving…" : "Save"}
				</Button>
			</div>

			<div className="flex-1 overflow-auto p-6">
				<div className="flex flex-col gap-3 max-w-3xl">
					{wf.steps.length === 0 && (
						<p className="text-sm text-muted-foreground italic py-2">
							No steps yet — add one below.
						</p>
					)}
					{wf.steps.map((step, si) => (
						// biome-ignore lint/suspicious/noArrayIndexKey: step index is stable within editor render
						<Card key={`editor-step-${si}`} className="border-border/60">
							<CardHeader className="bg-[hsl(240_10%_6%)] rounded-t-xl border-b border-border py-2.5 px-4">
								<div className="flex items-center gap-3">
									<span className="flex h-6 w-6 items-center justify-center rounded-full bg-muted text-xs font-bold text-muted-foreground shrink-0">
										{si + 1}
									</span>
									<Select
										value={step.kind}
										onValueChange={(k) => updStep(si, { kind: k })}
									>
										<SelectTrigger className="w-24 h-7">
											<SelectValue />
										</SelectTrigger>
										<SelectContent>
											<SelectItem value="sync">SYNC</SelectItem>
											<SelectItem value="fix">FIX</SelectItem>
											<SelectItem value="edit">EDIT</SelectItem>
										</SelectContent>
									</Select>
									<div className="flex-1" />
									<Button
										type="button"
										variant="ghost"
										size="icon"
										className="h-7 w-7"
										onClick={() =>
											mutSteps((ss) => ss.filter((_, j) => j !== si))
										}
									>
										<X className="h-3.5 w-3.5 text-muted-foreground" />
									</Button>
								</div>
							</CardHeader>

							<CardContent className="pt-3 pb-4 flex flex-col gap-3">
								<div>
									<p className="text-[10px] uppercase tracking-widest text-muted-foreground mb-1.5">
										Source
									</p>
									<EpEditRow
										ep={step.source}
										sources={sources}
										onChange={(p) => updSource(si, p)}
									/>
								</div>
								<div>
									<p className="text-[10px] uppercase tracking-widest text-muted-foreground mb-1.5">
										Target{step.targets.length > 1 ? "s" : ""}
									</p>
									<div className="flex flex-col gap-2">
										{step.targets.map((tgt, ti) => (
											<div
												key={`${tgt.source_id}-${tgt.resource}`}
												className="flex items-center gap-2"
											>
												<EpEditRow
													ep={tgt}
													sources={sources}
													onChange={(p) => updTarget(si, ti, p)}
												/>
												{step.targets.length > 1 && (
													<Button
														type="button"
														variant="ghost"
														size="icon"
														className="h-7 w-7 shrink-0"
														onClick={() => removeTarget(si, ti)}
													>
														<X className="h-3 w-3 text-muted-foreground" />
													</Button>
												)}
											</div>
										))}
										<button
											type="button"
											onClick={() => addTarget(si)}
											className="text-xs text-blue-400 hover:text-blue-300 text-left mt-0.5"
										>
											+ Add target
										</button>
									</div>
								</div>
								{si > 0 && (
									<div>
										<p className="text-[10px] uppercase tracking-widest text-muted-foreground mb-1.5">
											Run after (step IDs, comma-separated)
										</p>
										<Input
											className="h-7 text-xs font-mono"
											value={step.after?.join(", ") ?? ""}
											placeholder="Leave blank to run in parallel"
											onChange={(e) =>
												updStep(si, {
													after: e.target.value
														.split(",")
														.map((s) => s.trim())
														.filter(Boolean),
												})
											}
										/>
									</div>
								)}
							</CardContent>
						</Card>
					))}
					<button
						type="button"
						onClick={() => mutSteps((ss) => [...ss, blankStep(firstSrcId)])}
						className="w-full rounded-xl border border-dashed border-border py-3 text-sm text-muted-foreground hover:border-blue-700 hover:text-blue-400 transition-colors"
					>
						+ Add Step
					</button>
				</div>
			</div>
		</div>
	);
}
