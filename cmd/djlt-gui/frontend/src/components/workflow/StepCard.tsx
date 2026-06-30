import { CheckCircle, Clock, Pencil, Wrench, XCircle, Zap } from "lucide-react";
import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import type { Endpoint, Source, Step, StepDiff, StepResult } from "@/types";
import { TrackDiffTable } from "./TrackDiffTable";

// ── helpers ────────────────────────────────────────────────────────────────

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
		return <CheckCircle className="h-4 w-4 text-emerald-400" />;
	if (status === "failed") return <XCircle className="h-4 w-4 text-red-400" />;
	return <Clock className="h-4 w-4 text-purple-400" />;
}

// ── EndpointReadRow ────────────────────────────────────────────────────────
// Read-only mirror of EpEditRow — same proportions, same field order,
// non-interactive styling so it reads clearly without confusion with edit mode.

function EndpointReadRow({ ep, sources }: { ep: Endpoint; sources: Source[] }) {
	const sourceName =
		sources.find((s) => s.id === ep.source_id)?.name ??
		ep.source_id.slice(0, 8);

	return (
		<div className="flex gap-2 items-center">
			<div className="h-7 w-36 shrink-0 flex items-center px-2 rounded-md border border-border/40 bg-muted/20 text-xs text-muted-foreground truncate">
				{sourceName}
			</div>
			<div className="h-7 w-24 shrink-0 flex items-center px-2 rounded-md border border-border/40 bg-muted/20 text-xs text-muted-foreground">
				{ep.resource}
			</div>
			<div className="h-7 flex-1 min-w-0 flex items-center px-2 rounded-md border border-border/40 bg-muted/20 text-xs font-mono text-muted-foreground truncate">
				{ep.query ? (
					ep.query
				) : (
					<span className="text-muted-foreground/40">—</span>
				)}
			</div>
		</div>
	);
}

// ── StepCard ───────────────────────────────────────────────────────────────

interface StepCardProps {
	step: Step;
	index: number;
	sources: Source[];
	diff?: StepDiff;
	result?: StepResult;
	showResult: boolean;
}

export function StepCard({
	step,
	index,
	sources,
	diff,
	result,
	showResult,
}: StepCardProps) {
	const [showUnchanged, setShowUnchanged] = useState(true);
	const removedSet = new Set(diff?.removed.map((t) => t.id) ?? []);
	const unchanged = (diff?.current ?? []).filter((t) => !removedSet.has(t.id));

	const borderClass =
		result?.status === "success"
			? "border-emerald-900"
			: result?.status === "failed"
				? "border-red-900"
				: result?.status === "blocked"
					? "border-purple-900"
					: "border-border/60";

	const hasDiff = diff && step.kind === "sync";

	return (
		<Card className={borderClass}>
			{/* ── header ── */}
			<CardHeader className="bg-[hsl(240_10%_6%)] rounded-t-xl border-b border-border py-2.5 px-4">
				<div className="flex items-center gap-3">
					<span className="flex h-6 w-6 items-center justify-center rounded-full bg-muted text-xs font-bold text-muted-foreground shrink-0">
						{index + 1}
					</span>
					<Badge
						variant={kindVariant(step.kind)}
						className="flex items-center gap-1"
					>
						{kindIcon(step.kind)} {step.kind.toUpperCase()}
					</Badge>
					<div className="flex-1" />
					{showResult && result && (
						<div className="flex items-center gap-1.5 shrink-0">
							{statusIcon(result.status)}
							<span className="text-xs text-muted-foreground capitalize">
								{result.status}
							</span>
						</div>
					)}
				</div>
			</CardHeader>

			{/* ── body ── */}
			<CardContent className="pt-3 pb-4 flex flex-col gap-3">
				{/* Source */}
				<div>
					<p className="text-[10px] uppercase tracking-widest text-muted-foreground mb-1.5">
						Source
					</p>
					<EndpointReadRow ep={step.source} sources={sources} />
				</div>

				{/* Targets */}
				<div>
					<p className="text-[10px] uppercase tracking-widest text-muted-foreground mb-1.5">
						Target{step.targets.length !== 1 ? "s" : ""}
					</p>
					<div className="flex flex-col gap-2">
						{step.targets.map((tgt, ti) => (
							// biome-ignore lint/suspicious/noArrayIndexKey: target index is stable within a step card render
							<EndpointReadRow key={ti} ep={tgt} sources={sources} />
						))}
					</div>
				</div>

				{/* Diff */}
				{hasDiff && (
					<>
						<Separator className="opacity-50" />
						{diff.added.length === 0 && diff.removed.length === 0 ? (
							<div className="flex items-center gap-2 text-xs text-emerald-400 bg-emerald-950/40 rounded-md px-3 py-2">
								<CheckCircle className="h-3.5 w-3.5" /> Already up to date
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
						{result?.error && (
							<p className="text-xs text-destructive">✗ {result.error}</p>
						)}
					</>
				)}
			</CardContent>
		</Card>
	);
}
