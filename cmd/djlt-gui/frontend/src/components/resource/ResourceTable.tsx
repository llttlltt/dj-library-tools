import { forwardRef } from "react";
import { TableVirtuoso } from "react-virtuoso";
import { Badge } from "@/components/ui/badge";
import {
	TableBody,
	TableCell,
	TableHead,
	TableRow,
} from "@/components/ui/table";
import type { GroupRow, QueryResult, TrackRow } from "@/types";

interface ResourceTableProps {
	result: QueryResult;
}

export function ResourceTable({ result }: ResourceTableProps) {
	// biome-ignore lint/suspicious/noExplicitAny: data is union of two row types
	const data: any[] = result.kind === "groups" ? result.groups : result.tracks;

	return (
		<div className="flex-1 rounded-xl border border-border/80 overflow-hidden bg-background h-full">
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
	);
}
