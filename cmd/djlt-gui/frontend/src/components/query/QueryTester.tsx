import { Loader2 } from "lucide-react";
import { useEffect, useState } from "react";
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
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";
import type { GroupRow, QueryResult, Source, TrackRow } from "@/types";
import { ListSources, PreviewQuery } from "../../../wailsjs/go/gui/App";

interface QueryTesterProps {
	open: boolean;
	onClose: () => void;
	initialSourceID?: string;
	initialResource?: string;
	initialQuery?: string;
}

const asSources = (x: unknown) => (x ?? []) as Source[];
const asQueryResult = (x: unknown) => x as QueryResult;

export function QueryTester({
	open,
	onClose,
	initialSourceID,
	initialResource,
	initialQuery,
}: QueryTesterProps) {
	const [sources, setSources] = useState<Source[]>([]);
	const [sourceID, setSourceID] = useState(initialSourceID ?? "");
	const [resource, setResource] = useState(initialResource ?? "tracks");
	const [query, setQuery] = useState(initialQuery ?? "");
	const [result, setResult] = useState<QueryResult | null>(null);
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		ListSources()
			.then((s) => setSources(asSources(s)))
			.catch(() => {});
	}, []);

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
			<SheetContent>
				<SheetHeader>
					<SheetTitle>Query Tester</SheetTitle>
					<SheetDescription>
						Test any query before using it in a Step.
					</SheetDescription>
				</SheetHeader>
				<QueryTesterControls
					sources={sources}
					sourceID={sourceID}
					resource={resource}
					query={query}
					loading={loading}
					onSourceID={setSourceID}
					onResource={setResource}
					onQuery={setQuery}
					onTest={handleTest}
				/>
				<QueryTesterResults result={result} error={error} />
			</SheetContent>
		</Sheet>
	);
}

// ── Shared sub-components used by both QueryTester (sheet) and QueryTesterView ──

interface ControlsProps {
	sources: Source[];
	sourceID: string;
	resource: string;
	query: string;
	loading: boolean;
	onSourceID: (v: string) => void;
	onResource: (v: string) => void;
	onQuery: (v: string) => void;
	onTest: () => void;
}

export function QueryTesterControls({
	sources,
	sourceID,
	resource,
	query,
	loading,
	onSourceID,
	onResource,
	onQuery,
	onTest,
}: ControlsProps) {
	return (
		<div className="flex flex-col gap-3">
			<div className="flex flex-col gap-1">
				<span className="text-[10px] uppercase tracking-widest text-muted-foreground">
					Source
				</span>
				<Select value={sourceID} onValueChange={onSourceID}>
					<SelectTrigger className="h-8 text-sm">
						<SelectValue placeholder="Select a source…" />
					</SelectTrigger>
					<SelectContent>
						{sources.map((s) => (
							<SelectItem key={s.id} value={s.id}>
								{s.name}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
			</div>

			<div className="flex flex-col gap-1">
				<span className="text-[10px] uppercase tracking-widest text-muted-foreground">
					Resource
				</span>
				<Select value={resource} onValueChange={onResource}>
					<SelectTrigger className="h-8 text-sm">
						<SelectValue />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="tracks">tracks</SelectItem>
						<SelectItem value="playlists">playlists</SelectItem>
						<SelectItem value="folders">folders</SelectItem>
					</SelectContent>
				</Select>
			</div>

			<div className="flex flex-col gap-1">
				<span className="text-[10px] uppercase tracking-widest text-muted-foreground">
					Query
				</span>
				<Input
					className="h-8 text-sm font-mono"
					value={query}
					onChange={(e) => onQuery(e.target.value)}
					placeholder="beatgrids-count:1 && bpm:>120"
					onKeyDown={(e) => {
						if (e.key === "Enter" && sourceID && !loading) onTest();
					}}
				/>
			</div>

			<Button
				type="button"
				size="sm"
				onClick={onTest}
				disabled={loading || !sourceID}
				className="self-start"
			>
				{loading ? (
					<>
						<Loader2 className="h-3.5 w-3.5 mr-1.5 animate-spin" />
						Testing…
					</>
				) : (
					"Test"
				)}
			</Button>
		</div>
	);
}

interface ResultsProps {
	result: QueryResult | null;
	error: string;
}

export function QueryTesterResults({ result, error }: ResultsProps) {
	if (error) return <p className="text-sm text-destructive">{error}</p>;
	if (result === null) return null;

	const label =
		result.kind === "groups"
			? `Matched ${result.count} ${result.count !== 1 ? "items" : "item"}`
			: `Matched ${result.count} track${result.count !== 1 ? "s" : ""}`;

	const empty =
		result.kind === "groups"
			? "No playlists or folders matched."
			: "No tracks matched.";

	return (
		<div className="flex flex-col gap-2 flex-1 min-h-0">
			<Badge variant="secondary">{label}</Badge>

			{result.count === 0 ? (
				<p className="text-sm text-muted-foreground italic">{empty}</p>
			) : result.kind === "groups" ? (
				<div className="flex-1 overflow-auto rounded-md border border-border/60">
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead>Name</TableHead>
								<TableHead>Kind</TableHead>
								<TableHead>Parent</TableHead>
								<TableHead className="w-14 text-right">Items</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							{(result.groups as GroupRow[]).map((row) => (
								<TableRow key={row.id}>
									<TableCell className="text-sm font-medium truncate max-w-[160px]">
										{row.name || (
											<span className="text-muted-foreground italic">—</span>
										)}
									</TableCell>
									<TableCell className="text-sm text-muted-foreground">
										{row.kind}
									</TableCell>
									<TableCell className="text-sm text-muted-foreground truncate max-w-[120px]">
										{row.parent || "—"}
									</TableCell>
									<TableCell className="text-sm text-right font-mono text-muted-foreground">
										{row.items}
									</TableCell>
								</TableRow>
							))}
						</TableBody>
					</Table>
				</div>
			) : (
				<div className="flex-1 overflow-auto rounded-md border border-border/60">
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead>Title</TableHead>
								<TableHead>Artist</TableHead>
								<TableHead className="w-16 text-right">BPM</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							{(result.tracks as TrackRow[]).map((row) => (
								<TableRow key={row.id}>
									<TableCell className="text-sm truncate max-w-[180px]">
										{row.title || (
											<span className="text-muted-foreground italic">—</span>
										)}
									</TableCell>
									<TableCell className="text-sm text-muted-foreground truncate max-w-[120px]">
										{row.artist || "—"}
									</TableCell>
									<TableCell className="text-sm text-right font-mono text-muted-foreground">
										{row.bpm}
									</TableCell>
								</TableRow>
							))}
						</TableBody>
					</Table>
				</div>
			)}
		</div>
	);
}
