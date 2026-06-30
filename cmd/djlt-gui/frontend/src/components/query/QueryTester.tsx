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
	onApply?: (query: string) => void;
}

const asSources = (x: unknown) => (x ?? []) as Source[];
const asQueryResult = (x: unknown) => x as QueryResult;

export function QueryTester({
	open,
	onClose,
	initialSourceID,
	initialResource,
	initialQuery,
	onApply,
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
	onApply?: (query: string) => void;
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
	onApply,
}: ControlsProps) {
	const selectedSource = sources.find((s) => s.id === sourceID);
	const supportsFolders = selectedSource?.provider === "rb";

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
						{supportsFolders && (
							<SelectItem value="folders">folders</SelectItem>
						)}
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

			<div className="flex gap-2">
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

				{onApply && (
					<Button
						type="button"
						variant="outline"
						size="sm"
						onClick={() => onApply(query)}
						disabled={!query}
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

	if (result.count === 0) {
		return (
			<div className="flex flex-col gap-2 flex-1 min-h-0">
				<Badge variant="secondary" className="w-fit">
					{label}
				</Badge>
				<p className="text-sm text-muted-foreground italic">{empty}</p>
			</div>
		);
	}

	// biome-ignore lint/suspicious/noExplicitAny: data is union of two row types
	const data: any[] = result.kind === "groups" ? result.groups : result.tracks;

	return (
		<div className="flex flex-col gap-2 flex-1 min-h-0">
			<Badge variant="secondary" className="w-fit">
				{label}
			</Badge>
			<div className="flex-1 rounded-md border border-border/60 overflow-hidden bg-background">
				<TableVirtuoso
					data={data}
					totalCount={result.count}
					style={{ height: "100%" }}
					components={{
						Table: ({ ...props }) => (
							<Table {...props} className="border-collapse" />
						),
						TableHead: forwardRef((props, ref) => (
							<TableHeader
								{...props}
								ref={ref}
								className="sticky top-0 bg-background z-20"
							/>
						)),
						TableBody: forwardRef((props, ref) => (
							<TableBody {...props} ref={ref} />
						)),
						TableRow: (props) => <TableRow {...props} />,
					}}
					fixedHeaderContent={() => (
						<TableRow className="bg-background hover:bg-background">
							{result.kind === "groups" ? (
								<>
									<TableHead className="bg-background">Name</TableHead>
									<TableHead className="bg-background">Kind</TableHead>
									<TableHead className="bg-background">Parent</TableHead>
									<TableHead className="w-14 text-right bg-background">
										Items
									</TableHead>
								</>
							) : (
								<>
									<TableHead className="bg-background">Title</TableHead>
									<TableHead className="bg-background">Artist</TableHead>
									<TableHead className="w-16 text-right bg-background">
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
									<TableCell className="text-sm font-medium truncate max-w-[160px]">
										{g.name || (
											<span className="text-muted-foreground italic">—</span>
										)}
									</TableCell>
									<TableCell className="text-sm text-muted-foreground">
										{g.kind}
									</TableCell>
									<TableCell className="text-sm text-muted-foreground truncate max-w-[120px]">
										{g.parent || "—"}
									</TableCell>
									<TableCell className="text-sm text-right font-mono text-muted-foreground">
										{g.items}
									</TableCell>
								</>
							);
						}
						const t = row as TrackRow;
						return (
							<>
								<TableCell className="text-sm truncate max-w-[180px]">
									{t.title || (
										<span className="text-muted-foreground italic">—</span>
									)}
								</TableCell>
								<TableCell className="text-sm text-muted-foreground truncate max-w-[120px]">
									{t.artist || "—"}
								</TableCell>
								<TableCell className="text-sm text-right font-mono text-muted-foreground">
									{t.bpm}
								</TableCell>
							</>
						);
					}}
				/>
			</div>
		</div>
	);
}
