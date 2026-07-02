import {
	CheckCircle,
	Clock,
	FileText,
	FolderOpen,
	Pencil,
	Plus,
	Trash2,
	Wrench,
	X,
	XCircle,
	Zap,
} from "lucide-react";
import { useState } from "react";
import type { QueryTesterOpts } from "@/App";
import { EndpointEditor } from "@/components/endpoint/EndpointEditor";
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
import type {
	Connection,
	Endpoint,
	ProviderInfo,
	Step,
	StepDiff,
	StepResult,
} from "@/types";
import { OpenFileDialog } from "../../../wailsjs/go/gui/App";
import { TrackDiffTable } from "./TrackDiffTable";

// ── Helpers ────────────────────────────────────────────────────────────────

export function kindIcon(kind: string) {
	switch (kind.toLowerCase()) {
		case "sync":
			return <Zap className="h-3.5 w-3.5" />;
		case "fix":
			return <Wrench className="h-3.5 w-3.5" />;
		case "add":
			return <Plus className="h-3.5 w-3.5" />;
		case "remove":
			return <Trash2 className="h-3.5 w-3.5" />;
		case "m3u_export":
			return <FileText className="h-3.5 w-3.5" />;
		default:
			return <Pencil className="h-3.5 w-3.5" />;
	}
}

export function kindVariant(
	kind: string,
): "sync" | "fix" | "edit" | "add" | "remove" | "m3u_export" {
	if (kind === "sync") return "sync";
	if (kind === "fix") return "fix";
	if (kind === "add") return "add";
	if (kind === "remove") return "remove";
	if (kind === "m3u_export") return "m3u_export";
	return "edit";
}

export function statusIcon(status: string) {
	if (status === "success")
		return <CheckCircle className="h-4 w-4 text-emerald-500" />;
	if (status === "failed") return <XCircle className="h-4 w-4 text-rose-500" />;
	return <Clock className="h-4 w-4 text-purple-400" />;
}

// ── EndpointReadRow ────────────────────────────────────────────────────────

function EndpointReadRow({
	ep,
	connections,
}: {
	ep: Endpoint;
	connections: Connection[];
}) {
	let connectionName = ep.connection_id;
	if (ep.connection_id === "m3u") connectionName = "AD-HOC M3U";
	else if (ep.connection_id === "m3u8") connectionName = "AD-HOC M3U8";
	else {
		connectionName =
			connections.find((c) => c.id === ep.connection_id)?.name ??
			ep.connection_id.slice(0, 8);
	}

	return (
		<div className="flex flex-nowrap gap-2.5 items-center w-full min-w-0">
			<div className="h-9 w-40 shrink-0 flex items-center px-3 rounded-lg border border-border/60 bg-muted/30 text-sm font-medium text-foreground truncate">
				{connectionName}
			</div>
			<div className="h-9 w-40 shrink-0 flex items-center justify-center px-2 rounded-lg border border-border/60 bg-muted/30 text-xs font-mono font-medium text-muted-foreground truncate">
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

// ── TargetsSection ──────────────────────────────────────────────────────────

function TargetsSection({
	step,
	isEdit,
	connections,
	providers,
	updTarget,
	removeTarget,
	addTarget,
	onOpenQueryTester,
}: {
	step: Step;
	isEdit: boolean;
	connections: Connection[];
	providers: ProviderInfo[];
	updTarget: (ti: number, patch: Partial<Endpoint>) => void;
	removeTarget: (ti: number) => void;
	addTarget: () => void;
	onOpenQueryTester?: (opts?: QueryTesterOpts) => void;
}) {
	return (
		<div className="space-y-1.5">
			<p className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground pl-0.5">
				{step.kind === "add"
					? "Target Group Name (Query)"
					: step.kind === "remove"
						? "Target Group(s)"
						: `Target${step.targets.length !== 1 ? "s" : ""}`}
			</p>
			<div className="flex flex-col gap-2.5">
				{step.targets.map((tgt, ti) => (
					<div
						key={tgt.connection_id}
						className="flex flex-nowrap items-center gap-2.5 w-full min-w-0"
					>
						<div className="flex-1 min-w-0">
							{isEdit ? (
								<EndpointEditor
									endpoint={tgt}
									connections={connections}
									providers={providers}
									isTarget
									onChange={(p) => updTarget(ti, p)}
									onOpenQueryTester={onOpenQueryTester}
									layout="row"
								/>
							) : (
								<EndpointReadRow ep={tgt} connections={connections} />
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
	);
}

// ── M3uExportSection ────────────────────────────────────────────────────────

function M3uExportSection({
	step,
	index,
	isEdit,
	onChange,
}: {
	step: Step;
	index: number;
	isEdit: boolean;
	onChange?: (patch: Partial<Step>) => void;
}) {
	return (
		<div className="space-y-1.5 pt-1">
			<p className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground pl-0.5">
				Export Path
			</p>
			<div className="flex gap-2">
				<Input
					className="h-9 text-sm bg-background/50 font-mono"
					value={(step.options?.path as string) ?? ""}
					onChange={(e) =>
						onChange?.({
							options: { ...step.options, path: e.target.value },
						})
					}
					placeholder="/path/to/playlist.m3u"
					disabled={!isEdit}
				/>
				{isEdit && (
					<Button
						variant="outline"
						size="icon"
						className="h-9 w-9 shrink-0 bg-background/50"
						onClick={async () => {
							const path = await OpenFileDialog("");
							if (path) {
								onChange?.({ options: { ...step.options, path } });
							}
						}}
					>
						<FolderOpen className="h-4 w-4 text-muted-foreground" />
					</Button>
				)}
			</div>
			<div className="flex items-center gap-2 mt-2 px-0.5">
				<input
					type="checkbox"
					id={`append-${index}`}
					className="rounded border-border bg-background/50"
					checked={!!step.options?.append}
					disabled={!isEdit}
					onChange={(e) =>
						onChange?.({
							options: { ...step.options, append: e.target.checked },
						})
					}
				/>
				<label
					htmlFor={`append-${index}`}
					className="text-[11px] font-medium text-muted-foreground cursor-pointer select-none"
				>
					Append to existing file
				</label>
			</div>
		</div>
	);
}

// ── DiffSection ─────────────────────────────────────────────────────────────

function DiffSection({
	diffs,
	showUnchanged,
	onToggleUnchanged,
}: {
	diffs: StepDiff[];
	showUnchanged: boolean;
	onToggleUnchanged: () => void;
}) {
	return (
		<div className="space-y-4">
			<Separator className="opacity-40 my-1" />
			{diffs.map((diff) => {
				const removedSet = new Set(diff.removed.map((t) => t.id));
				const unchanged = diff.current.filter((t) => !removedSet.has(t.id));

				return (
					<div
						key={`${diff.step_id}-${diff.target_name}`}
						className="space-y-3"
					>
						{diff.added.length === 0 && diff.removed.length === 0 ? (
							<div className="flex items-center gap-2 text-xs font-medium text-emerald-400 bg-emerald-950/20 border border-emerald-500/20 rounded-xl px-3.5 py-2.5">
								<CheckCircle className="h-4 w-4 shrink-0 text-emerald-500" />{" "}
								{diff.target_name}: Already up to date
							</div>
						) : (
							<TrackDiffTable
								target={diff.target_name}
								added={diff.added}
								removed={diff.removed}
								unchanged={unchanged}
								showUnchanged={showUnchanged}
								onToggleUnchanged={onToggleUnchanged}
							/>
						)}
					</div>
				);
			})}
		</div>
	);
}

// ── StepCard ───────────────────────────────────────────────────────────────

interface StepCardProps {
	mode: "edit" | "view" | "applying";
	step: Step;
	index: number;
	connections: Connection[];
	providers?: ProviderInfo[];
	// Edit mode props
	onChange?: (patch: Partial<Step>) => void;
	onDelete?: () => void;
	onOpenQueryTester?: (opts?: QueryTesterOpts) => void;
	// View/Applying mode props
	diffs?: StepDiff[];
	result?: StepResult;
}

export function StepCard({
	mode,
	step,
	index,
	connections,
	providers = [],
	onChange,
	onDelete,
	onOpenQueryTester,
	diffs = [],
	result,
}: StepCardProps) {
	const [showUnchanged, setShowUnchanged] = useState(true);
	const isEdit = mode === "edit";
	const showResult = mode === "applying";

	// Logic for Diffs
	const hasDiffs =
		diffs.length > 0 &&
		(step.kind === "sync" ||
			step.kind === "add" ||
			step.kind === "remove" ||
			step.kind === "m3u_export");

	// Edit Handlers
	const updateConnection = (patch: Partial<Endpoint>) =>
		onChange?.({ source: { ...step.source, ...patch } });

	const updTarget = (ti: number, patch: Partial<Endpoint>) => {
		const tgts = [...step.targets];
		tgts[ti] = { ...tgts[ti], ...patch };
		onChange?.({ targets: tgts });
	};

	const addTarget = () => {
		const sortedConnections = [...connections].sort((a, b) =>
			a.name.localeCompare(b.name),
		);
		const firstTargetConnectionId =
			sortedConnections.find((c) => {
				const p = providers.find((prov) => prov.name === c.provider);
				return p?.capabilities.CanWrite ?? true;
			})?.id ??
			sortedConnections[0]?.id ??
			"";

		onChange?.({
			targets: [
				...step.targets,
				{
					connection_id: firstTargetConnectionId,
					resource: "playlists",
					query: "",
				},
			],
		});
	};

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
								<SelectItem value="add">ADD</SelectItem>
								<SelectItem value="remove">REMOVE</SelectItem>
								<SelectItem value="m3u_export">M3U EXPORT</SelectItem>
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
				{/* Connection Section */}
				<div className="space-y-1.5">
					<p className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground pl-0.5">
						{step.kind === "add"
							? "Source Selection"
							: step.kind === "remove"
								? "Tracks to Remove"
								: "Source Connection"}
					</p>
					{isEdit ? (
						<EndpointEditor
							endpoint={step.source}
							connections={connections}
							providers={providers}
							onChange={updateConnection}
							onOpenQueryTester={onOpenQueryTester}
							layout="row"
						/>
					) : (
						<EndpointReadRow ep={step.source} connections={connections} />
					)}
				</div>

				{/* M3U Export Path Section */}
				{step.kind === "m3u_export" && (
					<M3uExportSection
						step={step}
						index={index}
						isEdit={isEdit}
						onChange={onChange}
					/>
				)}

				{/* Targets Section */}
				{step.kind !== "m3u_export" && (
					<TargetsSection
						step={step}
						isEdit={isEdit}
						connections={connections}
						providers={providers}
						updTarget={updTarget}
						removeTarget={removeTarget}
						addTarget={addTarget}
						onOpenQueryTester={onOpenQueryTester}
					/>
				)}

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
										.map((c) => c.trim())
										.filter(Boolean),
								})
							}
						/>
					</div>
				)}

				{/* Diff & Results Section (View/Applying Mode) */}
				{!isEdit && hasDiffs && (
					<DiffSection
						diffs={diffs}
						showUnchanged={showUnchanged}
						onToggleUnchanged={() => setShowUnchanged((v) => !v)}
					/>
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
