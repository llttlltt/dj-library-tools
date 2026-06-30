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
import type { Source, TrackRow } from "@/types";
import { ListSources, PreviewQuery } from "../../../wailsjs/go/gui/App";

interface QueryTesterProps {
	open: boolean;
	onClose: () => void;
	initialSourceID?: string;
	initialResource?: string;
	initialQuery?: string;
}

const asSources = (x: unknown) => (x ?? []) as Source[];
const asRows = (x: unknown) => (x ?? []) as TrackRow[];

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
	const [results, setResults] = useState<TrackRow[] | null>(null);
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	// Load sources once on mount
	useEffect(() => {
		ListSources()
			.then((s) => setSources(asSources(s)))
			.catch(() => {});
	}, []);

	// Sync initial values when they change (e.g. opened from a different step)
	useEffect(() => {
		setSourceID(initialSourceID ?? "");
	}, [initialSourceID]);

	useEffect(() => {
		setResource(initialResource ?? "tracks");
	}, [initialResource]);

	useEffect(() => {
		setQuery(initialQuery ?? "");
	}, [initialQuery]);

	// Clear results whenever inputs change
	// biome-ignore lint/correctness/useExhaustiveDependencies: intentionally watching value changes to clear stale results
	useEffect(() => {
		setResults(null);
		setError("");
	}, [sourceID, resource, query]);

	async function handleTest() {
		setLoading(true);
		setError("");
		setResults(null);
		try {
			const rows = asRows(await PreviewQuery(sourceID, resource, query));
			setResults(rows);
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

				{/* ── Controls ── */}
				<div className="flex flex-col gap-3">
					{/* Source */}
					<div className="flex flex-col gap-1">
						<span className="text-[10px] uppercase tracking-widest text-muted-foreground">
							Source
						</span>
						<Select value={sourceID} onValueChange={setSourceID}>
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

					{/* Resource */}
					<div className="flex flex-col gap-1">
						<span className="text-[10px] uppercase tracking-widest text-muted-foreground">
							Resource
						</span>
						<Select value={resource} onValueChange={setResource}>
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

					{/* Query */}
					<div className="flex flex-col gap-1">
						<span className="text-[10px] uppercase tracking-widest text-muted-foreground">
							Query
						</span>
						<Input
							className="h-8 text-sm font-mono"
							value={query}
							onChange={(e) => setQuery(e.target.value)}
							placeholder="beatgrids-count:1 && bpm:>120"
							onKeyDown={(e) => {
								if (e.key === "Enter" && sourceID && !loading) handleTest();
							}}
						/>
					</div>

					<Button
						type="button"
						size="sm"
						onClick={handleTest}
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

				{/* ── Results ── */}
				{error && <p className="text-sm text-destructive">{error}</p>}

				{results !== null && !error && (
					<div className="flex flex-col gap-2 flex-1 min-h-0">
						<div className="flex items-center gap-2">
							<Badge variant="secondary">
								Matched {results.length} track{results.length !== 1 ? "s" : ""}
							</Badge>
						</div>

						{results.length === 0 ? (
							<p className="text-sm text-muted-foreground italic">
								No tracks matched.
							</p>
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
										{results.map((row) => (
											<TableRow key={row.id}>
												<TableCell className="text-sm truncate max-w-[180px]">
													{row.title || (
														<span className="text-muted-foreground italic">
															—
														</span>
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
				)}
			</SheetContent>
		</Sheet>
	);
}
