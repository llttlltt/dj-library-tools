import {
	CheckCircle,
	Clock,
	FlaskConical,
	Pencil,
	Wrench,
	X,
	XCircle,
	Zap,
} from "lucide-react";
import { useState } from "react";
import type { QueryTesterOpts } from "@/App";
import { Badge } from "@/components/ui/badge";
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
import type { Endpoint, Source, Step, StepDiff, StepResult } from "@/types";
import { TrackDiffTable } from "./TrackDiffTable";

// ── Helpers ────────────────────────────────────────────────────────────────

export function kindIcon(kind: string) {
	switch (kind.toLowerCase()) {
		case "sync":
			return <Zap className="h-3.5 w-3.5" />;
		case "fix":
			return <Wrench className="h-3.5 w-3.5" />;
		default:
			return <Pencil className="h-3.5 w-3.5" />;
	}
}

export function kindVariant(kind: string): "sync" | "fix" | "edit" {
	if (kind === "sync") return "sync";
	if (kind === "fix") return "fix";
	return "edit";
}

export function statusIcon(status: string) {
	if (status === "success")
		return <CheckCircle className="h-4 w-4 text-emerald-500" />;
	if (status === "failed") return <XCircle className="h-4 w-4 text-rose-500" />;
	return <Clock className="h-4 w-4 text-purple-400" />;
}

// ── EpEditRow ──────────────────────────────────────────────────────────────

interface EpEditRowProps {
	ep: Endpoint;
	sources: Source[];
	onChange: (p: Partial<Endpoint>) => void;
	onOpenQueryTester?: (opts?: QueryTesterOpts) => void;
}

export function EpEditRow({
	ep,
	sources,
	onChange,
	onOpenQueryTester,
}: EpEditRowProps) {
	return (
		<div className="flex flex-nowrap gap-2.5 items-center w-full min-w-0">
			<Select
				value={ep.source_id}
				onValueChange={(v) => onChange({ source_id: v })}
			>
				<SelectTrigger className="w-40 h-9 text-sm shrink-0 bg-background/50">
					<SelectValue placeholder="Source" />
				</SelectTrigger>
				<SelectContent>
					{[...sources]
						.sort((a, b) => a.name.localeCompare(b.name))
						.map((s) => (
							<SelectItem key={s.id} value={s.id}>
								{s.name}
							</SelectItem>
						))}
				</SelectContent>
			</Select>
			<Input
				className="h-9 text-sm w-24 shrink-0 bg-background/50"
				value={ep.resource}
				onChange={(e) => onChange({ resource: e.target.value })}
				placeholder="resource"
			/>
			<Input
				className="h-9 text-sm flex-1 min-w-0 bg-background/50"
				value={ep.query ?? ""}
				onChange={(e) => onChange({ query: e.target.value })}
				placeholder="query (optional)"
			/>
			{onOpenQueryTester ? (
				<Button
					type="button"
					variant="ghost"
					size="icon"
					className="h-8.5 w-8.5 shrink-0 hover:bg-secondary"
					title="Test query"
					onClick={() =>
						onOpenQueryTester({
							sourceID: ep.source_id,
							resource: ep.resource,
							query: ep.query ?? "",
							onApply: (q) => onChange({ query: q }),
						})
					}
				>
					<FlaskConical className="h-4 w-4 text-muted-foreground" />
				</Button>
			) : (
				// Matches the visual and layout spacing of the tester button
				<div className="w-8.5 shrink-0" />
			)}
		</div>
	);
}

// ── EndpointReadRow ────────────────────────────────────────────────────────

function EndpointReadRow({ ep, sources }: { ep: Endpoint; sources: Source[] }) {
	const sourceName =
		sources.find((s) => s.id === ep.source_id)?.name ??
		ep.source_id.slice(0, 8);

	return (
		<div className="flex flex-nowrap gap-2.5 items-center w-full min-w-0">
			<div className="h-9 w-40 shrink-0 flex items-center px-3 rounded-lg border border-border/60 bg-muted/30 text-sm font-medium text-foreground truncate">
				{sourceName}
			</div>
			<div className="h-9 w-24 shrink-0 flex items-center justify-center px-2 rounded-lg border border-border/60 bg-muted/30 text-xs font-mono font-medium text-muted-foreground">
				{ep.resource}
			</div>
			<div className="h-9 flex-1 min-w-0 flex items-center px-3 rounded-lg border border-border/60 bg-muted/20 text-sm font-mono text-muted-foreground truncate">
				{ep.query || (
					<span className="text-muted-foreground/30 font-sans italic">—</span>
				)}
			</div>
			{/* Reserved layout space matching the action action trigger width */}
			<div className="w-8.5 shrink-0" />
		</div>
	);
}

// ── StepCard ───────────────────────────────────────────────────────────────

interface StepCardProps {
	mode: "edit" | "view" | "applying";
	step: Step;
	index: number;
	sources: Source[];
	// Edit mode props
	onChange?: (patch: Partial<Step>) => void;
	onDelete?: () => void;
	onOpenQueryTester?: (opts?: QueryTesterOpts) => void;
	// View/Applying mode props
	diff?: StepDiff;
	result?: StepResult;
}

export function StepCard({
	mode,
	step,
	index,
	sources,
	onChange,
	onDelete,
	onOpenQueryTester,
	diff,
	result,
}: StepCardProps) {
	const [showUnchanged, setShowUnchanged] = useState(true);
	const isEdit = mode === "edit";
	const showResult = mode === "applying";

	// Logic for Diffs
	const removedSet = new Set(diff?.removed.map((t) => t.id) ?? []);
	const unchanged = (diff?.current ?? []).filter((t) => !removedSet.has(t.id));
	const hasDiff = diff && step.kind === "sync";

	// Edit Handlers
	const updSource = (patch: Partial<Endpoint>) =>
		onChange?.({ source: { ...step.source, ...patch } });

	const updTarget = (ti: number, patch: Partial<Endpoint>) => {
		const tgts = [...step.targets];
		tgts[ti] = { ...tgts[ti], ...patch };
		onChange?.({ targets: tgts });
	};

	const addTarget = () =>
		onChange?.({
			targets: [
				...step.targets,
				{ source_id: sources[0]?.id ?? "", resource: "playlists", query: "" },
			],
		});

	const removeTarget = (ti: number) =>
		onChange?.({ targets: step.targets.filter((_, j) => j !== ti) });

	// Styling
	const borderClass =
		showResult && result?.status === "success"
			? "border-emerald-500/35 bg-emerald-950/5"
			: showResult && result?.status === "failed"
				? "border-rose-500/35 bg-rose-950/5"
				: isEdit
					? "border-border/60"
					: "border-border/65";

	return (
		<Card
			className={`overflow-hidden transition-all duration-200 ${borderClass}`}
		>
			{/* ── Header ── Fixed height of [53px] prevents layout shifts between selector and badge */}
			<CardHeader
				className={`${isEdit ? "bg-[hsl(240_10%_6%)]" : "bg-secondary/25"} border-b border-border/80 h-[53px] py-0 px-4 flex flex-row items-center justify-between space-y-0`}
			>
				<div className="flex items-center gap-3">
					<span className="flex h-5 w-5 items-center justify-center rounded-full bg-secondary/85 text-[10px] font-bold text-muted-foreground shrink-0 border border-border/40">
						{index + 1}
					</span>

					{isEdit ? (
						<Select
							value={step.kind}
							onValueChange={(k) => onChange?.({ kind: k as Step["kind"] })}
						>
							<SelectTrigger className="w-28 h-8 text-xs font-semibold bg-background/50">
								<SelectValue />
							</SelectTrigger>
							<SelectContent>
								<SelectItem value="sync">SYNC</SelectItem>
								<SelectItem value="fix">FIX</SelectItem>
								<SelectItem value="edit">EDIT</SelectItem>
							</SelectContent>
						</Select>
					) : (
						<Badge
							variant={kindVariant(step.kind)}
							className="flex items-center gap-1 shrink-0 px-2 py-0.5"
						>
							{kindIcon(step.kind)} {step.kind.toUpperCase()}
						</Badge>
					)}
				</div>

				<div className="flex items-center gap-2">
					{showResult && result && (
						<div className="flex items-center gap-1.5 shrink-0">
							{statusIcon(result.status)}
							<span className="text-xs font-semibold text-muted-foreground capitalize">
								{result.status}
							</span>
						</div>
					)}
					{isEdit && (
						<Button
							type="button"
							variant="ghost"
							size="icon"
							className="h-8 w-8 hover:bg-secondary rounded-lg"
							onClick={onDelete}
						>
							<X className="h-4 w-4 text-muted-foreground" />
						</Button>
					)}
				</div>
			</CardHeader>

			{/* ── Body ── */}
			<CardContent className="pt-4 pb-5 flex flex-col gap-4">
				{/* Source Section */}
				<div className="space-y-1.5">
					<p className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground pl-0.5">
						Source
					</p>
					{isEdit ? (
						<EpEditRow
							ep={step.source}
							sources={sources}
							onChange={updSource}
							onOpenQueryTester={onOpenQueryTester}
						/>
					) : (
						<EndpointReadRow ep={step.source} sources={sources} />
					)}
				</div>

				{/* Targets Section */}
				<div className="space-y-1.5">
					<p className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground pl-0.5">
						Target{step.targets.length !== 1 ? "s" : ""}
					</p>
					<div className="flex flex-col gap-2.5">
						{step.targets.map((tgt, ti) => (
							<div
								key={tgt.source_id}
								className="flex flex-nowrap items-center gap-2.5 w-full min-w-0"
							>
								<div className="flex-1 min-w-0">
									{isEdit ? (
										<EpEditRow
											ep={tgt}
											sources={sources}
											onChange={(p) => updTarget(ti, p)}
											onOpenQueryTester={onOpenQueryTester}
										/>
									) : (
										<EndpointReadRow ep={tgt} sources={sources} />
									)}
								</div>
								{isEdit && step.targets.length > 1 && (
									<Button
										type="button"
										variant="ghost"
										size="icon"
										className="h-8.5 w-8.5 shrink-0 hover:bg-secondary rounded-lg"
										onClick={() => removeTarget(ti)}
									>
										<X className="h-4 w-4 text-muted-foreground" />
									</Button>
								)}
							</div>
						))}
						{isEdit && (
							<button
								type="button"
								onClick={addTarget}
								className="text-xs text-blue-400 hover:text-blue-300 transition-colors text-left font-medium mt-1 pl-0.5"
							>
								+ Add target
							</button>
						)}
					</div>
				</div>

				{/* Run After Section (Edit Mode only) */}
				{isEdit && index > 0 && (
					<div className="space-y-1.5 pt-1">
						<p className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground pl-0.5">
							Run after (step IDs, comma-separated)
						</p>
						<Input
							className="h-9 text-sm font-mono bg-background/50"
							value={step.after?.join(", ") ?? ""}
							placeholder="Leave blank to run in parallel"
							onChange={(e) =>
								onChange?.({
									after: e.target.value
										.split(",")
										.map((s) => s.trim())
										.filter(Boolean),
								})
							}
						/>
					</div>
				)}

				{/* Diff & Results Section (View/Applying Mode) */}
				{!isEdit && hasDiff && (
					<>
						<Separator className="opacity-40 my-1" />
						{diff.added.length === 0 && diff.removed.length === 0 ? (
							<div className="flex items-center gap-2 text-xs font-medium text-emerald-400 bg-emerald-950/20 border border-emerald-500/20 rounded-xl px-3.5 py-2.5">
								<CheckCircle className="h-4 w-4 shrink-0 text-emerald-500" />{" "}
								Already up to date
							</div>
						) : (
							<TrackDiffTable
								target={diff.target_name}
								added={diff.added}
								removed={diff.removed}
								unchanged={unchanged}
								showUnchanged={showUnchanged}
								onToggleUnchanged={() => setShowUnchanged((v) => !v)}
							/>
						)}
					</>
				)}

				{result?.error && (
					<div className="p-3.5 rounded-xl border border-destructive/20 bg-destructive/5 text-xs text-destructive font-mono leading-relaxed mt-1">
						<span className="font-semibold mr-1">Error:</span> {result.error}
					</div>
				)}
			</CardContent>
		</Card>
	);
}
