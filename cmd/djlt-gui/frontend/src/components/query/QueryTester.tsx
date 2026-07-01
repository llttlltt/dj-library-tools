import { useAtom } from "@effect-atom/atom-react";
import { Loader2 } from "lucide-react";
import { forwardRef, useEffect, useState } from "react";
import { TableVirtuoso } from "react-virtuoso";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from "@/components/ui/sheet";
import {
	TableBody,
	TableCell,
	TableHead,
	TableRow,
} from "@/components/ui/table";
import { runtime } from "@/lib/runtime";
import { cn } from "@/lib/utils";
import { loadProviders, providersAtom } from "@/store/providers";
import { loadSources, sourcesAtom } from "@/store/sources";
import type {
	GroupRow,
	ProviderInfo,
	QueryResult,
	Source,
	TrackRow,
} from "@/types";
import { PreviewQuery } from "../../../wailsjs/go/gui/App";

interface QueryTesterProps {
	open: boolean;
	onClose: () => void;
	initialSourceID?: string;
	initialResource?: string;
	initialQuery?: string;
	isTarget?: boolean;
	onApply?: (query: string) => void;
}

const asQueryResult = (x: unknown) => x as QueryResult;

export function QueryTester({
	open,
	onClose,
	initialSourceID,
	initialResource,
	initialQuery,
	isTarget,
	onApply,
}: QueryTesterProps) {
	const [sources] = useAtom(sourcesAtom);
	const [providers] = useAtom(providersAtom);

	const [sourceID, setSourceID] = useState(initialSourceID ?? "");
	const [resource, setResource] = useState(initialResource ?? "tracks");
	const [query, setQuery] = useState(initialQuery ?? "");
	const [result, setResult] = useState<QueryResult | null>(null);
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		runtime.runPromise(loadSources);
		runtime.runPromise(loadProviders);
	}, []);

	// Automatically fix blank or invalid resource selections when source or providers change
	useEffect(() => {
		if (sourceID && providers.length > 0) {
			const selectedSource = sources.find((s) => s.id === sourceID);
			const provider = providers.find(
				(p) => p.name === selectedSource?.provider,
			);
			const availableResources = (provider?.resources ?? []).filter((r) => {
				if (!isTarget) return true;
				return r.can_write;
			});

			if (availableResources.length > 0) {
				const isValid = availableResources.some((r) => r.name === resource);
				if (!isValid) {
					setResource(availableResources[0].name);
				}
			}
		}
	}, [sourceID, providers, resource, isTarget, sources]);

	useEffect(() => {
		setSourceID(initialSourceID ?? "");
	}, [initialSourceID]);

	useEffect(() => {
		setResource(initialResource ?? "tracks");
	}, [initialResource]);

	useEffect(() => {
		setQuery(initialQuery ?? "");
	}, [initialQuery]);

	// biome-ignore lint/correctness/useExhaustiveDependencies: intentionally watching value changes to clear stale results
	useEffect(() => {
		setResult(null);
		setError("");
	}, [sourceID, resource, query]);

	async function handleTest() {
		setLoading(true);
		setError("");
		setResult(null);
		try {
			setResult(asQueryResult(await PreviewQuery(sourceID, resource, query)));
		} catch (e) {
			setError(String(e));
		}
		setLoading(false);
	}

	return (
		<Sheet open={open} onOpenChange={(o) => !o && onClose()}>
			<SheetContent className="flex flex-col h-full sm:max-w-xl md:max-w-2xl">
				<SheetHeader className="shrink-0 mb-4">
					<SheetTitle>Query Tester</SheetTitle>
					<SheetDescription>
						Test any query before using it in a Step.
					</SheetDescription>
				</SheetHeader>

				<div className="flex flex-col gap-6 flex-1 min-h-0">
					<QueryTesterControls
						sources={sources}
						providers={providers}
						sourceID={sourceID}
						resource={resource}
						query={query}
						loading={loading}
						isTarget={isTarget}
						onSourceID={setSourceID}
						onResource={setResource}
						onQuery={setQuery}
						onTest={handleTest}
						onApply={
							onApply
								? (q) => {
										onApply(q);
										onClose();
									}
								: undefined
						}
					/>
					<QueryTesterResults result={result} error={error} />
				</div>
			</SheetContent>
		</Sheet>
	);
}

// ── Shared sub-components used by both QueryTester (sheet) and QueryTesterView ──

interface ControlsProps {
	sources: Source[];
	providers: ProviderInfo[];
	sourceID: string;
	resource: string;
	query: string;
	loading: boolean;
	isTarget?: boolean;
	onSourceID: (v: string) => void;
	onResource: (v: string) => void;
	onQuery: (v: string) => void;
	onTest: () => void;
	onApply?: (query: string) => void;
}

export function QueryTesterControls({
	sources,
	providers,
	sourceID,
	resource,
	query,
	loading,
	isTarget,
	onSourceID,
	onResource,
	onQuery,
	onTest,
	onApply,
}: ControlsProps) {
	const handleSourceChange = (v: string) => {
		const newSrc = sources.find((s) => s.id === v);
		const provider = providers.find((p) => p.name === newSrc?.provider);
		const availableResources = (provider?.resources ?? []).filter((r) => {
			if (!isTarget) return true;
			return r.can_write;
		});
		if (availableResources.length > 0) {
			const isValid = availableResources.some((r) => r.name === resource);
			if (!isValid) {
				onResource(availableResources[0].name);
			}
		}
		onSourceID(v);
	};

	const selectedSource = sources.find((s) => s.id === sourceID);
	const provider = providers.find((p) => p.name === selectedSource?.provider);
	const availableResources = (provider?.resources ?? []).filter((r) => {
		if (!isTarget) return true;
		return r.can_write;
	});

	// If the current resource isn't writable but we're in target mode, this is an invalid state
	const currentRes = provider?.resources.find((r) => r.name === resource);
	const isInvalidTarget = isTarget && currentRes && !currentRes.can_write;

	return (
		<div className="space-y-4 bg-secondary/20 p-4 rounded-xl border border-border/40 shrink-0">
			{/* Input Fields Grid */}
			<div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
				<div className="flex flex-col gap-1.5">
					<span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
						Source
					</span>
					<Select value={sourceID} onValueChange={handleSourceChange}>
						<SelectTrigger className="h-8.5 text-sm bg-background/50">
							<SelectValue placeholder="Select a source…" />
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
				</div>

				<div className="flex flex-col gap-1.5">
					<span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
						Resource
					</span>
					<Select value={resource} onValueChange={onResource}>
						<SelectTrigger
							className={cn(
								"h-8.5 text-sm bg-background/50 transition-colors",
								isInvalidTarget &&
									"border-destructive/60 bg-destructive/5 text-destructive",
							)}
						>
							<SelectValue />
						</SelectTrigger>
						<SelectContent>
							{availableResources.map((r) => (
								<SelectItem key={r.name} value={r.name}>
									{r.name}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
					{isInvalidTarget && (
						<span className="text-[10px] text-destructive/80 font-medium animate-in fade-in slide-in-from-top-1">
							This resource is read-only and cannot be used as a target.
						</span>
					)}
				</div>
			</div>

			<div className="flex flex-col gap-1.5">
				<span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
					Query Expression
				</span>
				<Input
					className="h-8.5 text-sm font-mono bg-background/50"
					value={query}
					onChange={(e) => onQuery(e.target.value)}
					placeholder="beatgrids-count:1 && bpm:>120"
					onKeyDown={(e) => {
						if (e.key === "Enter" && sourceID && !loading) onTest();
					}}
				/>
			</div>

			{/* Actions Row */}
			<div className="flex gap-2 pt-1 border-t border-border/20">
				<Button
					type="button"
					size="sm"
					onClick={onTest}
					disabled={loading || !sourceID}
					className="min-w-[80px]"
				>
					{loading ? (
						<>
							<Loader2 className="h-3.5 w-3.5 mr-1.5 animate-spin" />
							Testing…
						</>
					) : (
						"Test Query"
					)}
				</Button>

				{onApply && (
					<Button
						type="button"
						variant="secondary"
						size="sm"
						onClick={() => onApply(query)}
						disabled={!query || isInvalidTarget}
					>
						Use this query
					</Button>
				)}
			</div>
		</div>
	);
}

interface ResultsProps {
	result: QueryResult | null;
	error: string;
}

export function QueryTesterResults({ result, error }: ResultsProps) {
	if (error) {
		return (
			<div className="p-4 rounded-xl border border-destructive/20 bg-destructive/5 text-sm text-destructive font-mono overflow-auto shrink-0 leading-relaxed">
				<div className="font-semibold mb-1">Execution Error</div>
				{error}
			</div>
		);
	}
	if (result === null) return null;

	const label =
		result.kind === "groups"
			? `Matched ${result.count.toLocaleString()} ${result.count !== 1 ? "items" : "item"}`
			: `Matched ${result.count.toLocaleString()} track${result.count !== 1 ? "s" : ""}`;

	const empty =
		result.kind === "groups"
			? "No playlists or folders matched."
			: "No tracks matched.";

	if (result.count === 0) {
		return (
			<div className="flex flex-col items-center justify-center py-10 px-4 text-center border border-dashed border-border/60 rounded-xl bg-secondary/5 flex-1 min-h-0 space-y-2">
				<Badge
					variant="outline"
					className="text-muted-foreground py-0.5 px-2.5"
				>
					{label}
				</Badge>
				<p className="text-sm font-medium text-muted-foreground">{empty}</p>
				<p className="text-xs text-muted-foreground/60 max-w-xs">
					Try modifying the boolean keywords or parameters above and test again.
				</p>
			</div>
		);
	}

	// biome-ignore lint/suspicious/noExplicitAny: data is union of two row types
	const data: any[] = result.kind === "groups" ? result.groups : result.tracks;

	return (
		<div className="flex flex-col gap-3 flex-1 min-h-0">
			<div className="flex items-center justify-between">
				<span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
					Result Preview
				</span>
				<Badge
					variant="secondary"
					className="bg-emerald-500/10 text-emerald-500 border-emerald-500/20 py-0.5 px-2 text-xs font-medium"
				>
					{label}
				</Badge>
			</div>

			<div className="flex-1 rounded-xl border border-border/80 overflow-hidden bg-background">
				<TableVirtuoso
					data={data}
					totalCount={result.count}
					style={{ height: "100%" }}
					components={{
						Table: ({ ...props }) => (
							<table
								{...props}
								className="w-full border-collapse text-left text-sm"
							/>
						),
						TableHead: forwardRef((props, ref) => (
							<thead {...props} ref={ref} className="z-20" />
						)),
						TableBody: forwardRef((props, ref) => (
							<TableBody {...props} ref={ref} />
						)),
						TableRow: (props) => (
							<TableRow
								{...props}
								className="hover:bg-muted/30 transition-colors"
							/>
						),
					}}
					fixedHeaderContent={() => (
						<TableRow className="bg-secondary/40 border-b border-border/80">
							{result.kind === "groups" ? (
								<>
									<TableHead className="sticky top-0 bg-secondary/40 shadow-[0_1px_0_0_hsl(var(--border))] font-semibold text-xs py-2.5">
										Name
									</TableHead>
									<TableHead className="sticky top-0 bg-secondary/40 shadow-[0_1px_0_0_hsl(var(--border))] font-semibold text-xs py-2.5">
										Kind
									</TableHead>
									<TableHead className="sticky top-0 bg-secondary/40 shadow-[0_1px_0_0_hsl(var(--border))] font-semibold text-xs py-2.5">
										Parent
									</TableHead>
									<TableHead className="w-20 text-right sticky top-0 bg-secondary/40 shadow-[0_1px_0_0_hsl(var(--border))] font-semibold text-xs py-2.5">
										Items
									</TableHead>
								</>
							) : (
								<>
									<TableHead className="sticky top-0 bg-secondary/40 shadow-[0_1px_0_0_hsl(var(--border))] font-semibold text-xs py-2.5">
										Title
									</TableHead>
									<TableHead className="sticky top-0 bg-secondary/40 shadow-[0_1px_0_0_hsl(var(--border))] font-semibold text-xs py-2.5">
										Artist
									</TableHead>
									<TableHead className="w-20 text-right sticky top-0 bg-secondary/40 shadow-[0_1px_0_0_hsl(var(--border))] font-semibold text-xs py-2.5">
										BPM
									</TableHead>
								</>
							)}
						</TableRow>
					)}
					// biome-ignore lint/suspicious/noExplicitAny: row is union of TrackRow and GroupRow
					itemContent={(_, row: any) => {
						if (result.kind === "groups") {
							const g = row as GroupRow;
							return (
								<>
									<TableCell className="text-sm font-medium py-2 truncate max-w-[160px]">
										{g.name || (
											<span className="text-muted-foreground italic">—</span>
										)}
									</TableCell>
									<TableCell className="text-sm text-muted-foreground py-2">
										<Badge
											variant="outline"
											className="text-[10px] font-normal uppercase py-0 px-1.5 border-border/80 bg-background"
										>
											{g.kind}
										</Badge>
									</TableCell>
									<TableCell className="text-sm text-muted-foreground py-2 truncate max-w-[120px]">
										{g.parent || "—"}
									</TableCell>
									<TableCell className="text-sm text-right font-mono text-muted-foreground py-2 pr-4">
										{g.items}
									</TableCell>
								</>
							);
						}
						const t = row as TrackRow;
						const parsedBpm =
							typeof t.bpm === "string" ? Number.parseFloat(t.bpm) : t.bpm;
						return (
							<>
								<TableCell className="text-sm font-medium py-2 truncate max-w-[180px]">
									{t.title || (
										<span className="text-muted-foreground italic">—</span>
									)}
								</TableCell>
								<TableCell className="text-sm text-muted-foreground py-2 truncate max-w-[120px]">
									{t.artist || "—"}
								</TableCell>
								<TableCell className="text-sm text-right font-mono text-muted-foreground py-2 pr-4">
									{parsedBpm ? Math.round(parsedBpm) : "—"}
								</TableCell>
							</>
						);
					}}
				/>
			</div>
		</div>
	);
}
