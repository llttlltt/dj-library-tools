import {
	CheckCircle,
	ChevronRight,
	Clock,
	Pencil,
	Wrench,
	XCircle,
	Zap,
} from "lucide-react";
import React, { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import type { Endpoint, Step, StepDiff, StepResult } from "@/types";
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

function EndpointChip({ ep }: { ep: Endpoint }) {
	return (
		<span className="truncate text-xs font-mono text-muted-foreground bg-muted/40 px-1.5 py-0.5 rounded">
			{ep.resource}
			{ep.query ? ` · ${ep.query}` : ""}
		</span>
	);
}

// ── StepCard ───────────────────────────────────────────────────────────────

interface StepCardProps {
	step: Step;
	index: number;
	diff?: StepDiff;
	result?: StepResult;
	showResult: boolean;
}

export function StepCard({
	step,
	index,
	diff,
	result,
	showResult,
}: StepCardProps) {
	const [showUnchanged, setShowUnchanged] = useState(true);
	const removedSet = new Set(diff?.removed.map((t) => t.id) ?? []);
	const unchanged = (diff?.current ?? []).filter((t) => !removedSet.has(t.id));

	return (
		<Card
			className={
				result?.status === "success"
					? "border-emerald-900"
					: result?.status === "failed"
						? "border-red-900"
						: result?.status === "blocked"
							? "border-purple-900"
							: ""
			}
		>
			<CardHeader className="bg-[hsl(240_10%_6%)] rounded-t-xl border-b border-border py-3 px-4">
				<div className="flex flex-wrap items-center gap-3">
					<span className="flex h-6 w-6 items-center justify-center rounded-full bg-muted text-xs font-bold text-muted-foreground shrink-0">
						{index + 1}
					</span>
					<Badge
						variant={kindVariant(step.kind)}
						className="flex items-center gap-1"
					>
						{kindIcon(step.kind)} {step.kind.toUpperCase()}
					</Badge>
					<div className="flex items-center gap-1.5 min-w-0 flex-1 text-sm flex-wrap">
						<EndpointChip ep={step.source} />
						{step.targets.map((tgt, ti) => (
							// biome-ignore lint/suspicious/noArrayIndexKey: target index is stable within a step card render
							<React.Fragment key={`tgt-${ti}`}>
								<ChevronRight className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
								<EndpointChip ep={tgt} />
							</React.Fragment>
						))}
					</div>
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

			{diff && step.kind === "sync" && (
				<CardContent className="pt-3 pb-3">
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
						<p className="text-xs text-destructive mt-2">✗ {result.error}</p>
					)}
				</CardContent>
			)}
		</Card>
	);
}
