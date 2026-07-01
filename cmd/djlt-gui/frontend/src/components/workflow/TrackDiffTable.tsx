import type { TrackRow } from "@/types";

interface Props {
	target: string;
	added: TrackRow[];
	removed: TrackRow[];
	unchanged: TrackRow[];
	showUnchanged: boolean;
	onToggleUnchanged: () => void;
}

type RowKind = "add" | "remove" | "unchanged";
type DisplayRow = { kind: RowKind; track: TrackRow };

export function TrackDiffTable({
	target,
	added,
	removed,
	unchanged,
	showUnchanged,
	onToggleUnchanged,
}: Props) {
	const rows: DisplayRow[] = [
		...added.map((t) => ({ kind: "add" as RowKind, track: t })),
		...removed.map((t) => ({ kind: "remove" as RowKind, track: t })),
		...(showUnchanged
			? unchanged.map((t) => ({ kind: "unchanged" as RowKind, track: t }))
			: []),
	];

	return (
		<div className="space-y-2">
			{target && (
				<p className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground/80 pl-0.5">
					Target Partition:{" "}
					<span className="font-mono text-foreground lowercase">{target}</span>
				</p>
			)}
			<div className="rounded-xl border border-border/80 overflow-hidden bg-background">
				<table className="w-full border-collapse text-xs">
					<thead>
						<tr className="border-b border-border bg-secondary/40">
							<th className="w-8 py-2 pl-3" />
							<th className="py-2 px-2 text-left font-semibold text-muted-foreground">
								Title
							</th>
							<th className="py-2 px-2 text-left font-semibold text-muted-foreground hidden sm:table-cell">
								Artist
							</th>
							<th className="py-2 pr-4 text-right font-semibold text-muted-foreground w-16">
								BPM
							</th>
						</tr>
					</thead>
					<tbody className="divide-y divide-border/40">
						{rows.map(({ kind, track }) => {
							const parsedBpm =
								typeof track.bpm === "string"
									? Number.parseFloat(track.bpm)
									: track.bpm;
							return (
								<tr
									key={`${track.id}-${kind}`}
									className={`transition-colors ${
										kind === "add"
											? "bg-emerald-500/[0.03] hover:bg-emerald-500/[0.06]"
											: kind === "remove"
												? "bg-rose-500/[0.03] hover:bg-rose-500/[0.06]"
												: "opacity-60 hover:opacity-100 hover:bg-muted/30"
									}`}
								>
									<td className="py-2 pl-3 text-center">
										{kind === "add" ? (
											<span className="flex h-5 w-5 items-center justify-center rounded-md border border-emerald-500/20 bg-emerald-500/10 text-xs font-bold text-emerald-400">
												+
											</span>
										) : kind === "remove" ? (
											<span className="flex h-5 w-5 items-center justify-center rounded-md border border-rose-500/20 bg-rose-500/10 text-xs font-bold text-rose-450">
												−
											</span>
										) : (
											<span className="flex h-5 w-5 items-center justify-center rounded-md border border-border bg-muted/50 text-[10px] text-muted-foreground">
												·
											</span>
										)}
									</td>
									<td className="py-2 px-2 max-w-0 font-medium">
										<span className="truncate block">
											{track.title || track.id}
										</span>
									</td>
									<td className="py-2 px-2 text-muted-foreground/80 hidden sm:table-cell max-w-0">
										<span className="truncate block">
											{track.artist || "—"}
										</span>
									</td>
									<td className="py-2 pr-4 text-right text-muted-foreground font-mono">
										{parsedBpm ? Math.round(parsedBpm) : "—"}
									</td>
								</tr>
							);
						})}
					</tbody>
				</table>
			</div>
			{unchanged.length > 0 && (
				<button
					type="button"
					onClick={onToggleUnchanged}
					className="text-xs text-muted-foreground/80 hover:text-blue-400 transition-colors pl-1 font-medium flex items-center gap-1 mt-1"
				>
					{showUnchanged
						? `Hide ${unchanged.length.toLocaleString()} unchanged tracks`
						: `Show ${unchanged.length.toLocaleString()} unchanged tracks`}
				</button>
			)}
		</div>
	);
}
