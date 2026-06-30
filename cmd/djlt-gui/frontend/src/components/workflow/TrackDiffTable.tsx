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
		<div>
			{target && (
				<p className="text-xs text-muted-foreground mb-2 font-medium">
					{target}
				</p>
			)}
			<div className="rounded-md border border-border overflow-hidden">
				<table className="w-full text-xs">
					<thead>
						<tr className="border-b border-border bg-muted/20">
							<th className="w-5 py-1.5 pl-2" />
							<th className="py-1.5 px-2 text-left font-medium text-muted-foreground">
								Title
							</th>
							<th className="py-1.5 px-2 text-left font-medium text-muted-foreground hidden sm:table-cell">
								Artist
							</th>
							<th className="py-1.5 px-2 text-right font-medium text-muted-foreground w-12">
								BPM
							</th>
						</tr>
					</thead>
					<tbody>
						{rows.map(({ kind, track }) => (
							<tr
								key={`${track.id}-${kind}`}
								className={
									kind === "add"
										? "border-l-2 border-l-emerald-500 bg-emerald-950/30"
										: kind === "remove"
											? "border-l-2 border-l-red-500 bg-red-950/30"
											: "border-l-2 border-l-transparent opacity-50"
								}
							>
								<td className="pl-2 text-center font-mono font-bold">
									{kind === "add" ? (
										<span className="text-emerald-400">+</span>
									) : kind === "remove" ? (
										<span className="text-red-400">−</span>
									) : (
										<span className="text-muted-foreground">·</span>
									)}
								</td>
								<td className="py-1.5 px-2 max-w-0">
									<span className="truncate block">
										{track.title || track.id}
									</span>
								</td>
								<td className="py-1.5 px-2 text-muted-foreground hidden sm:table-cell max-w-0">
									<span className="truncate block">{track.artist}</span>
								</td>
								<td className="py-1.5 px-2 text-right text-muted-foreground font-mono">
									{track.bpm}
								</td>
							</tr>
						))}
					</tbody>
				</table>
			</div>
			{unchanged.length > 0 && (
				<button
					type="button"
					onClick={onToggleUnchanged}
					className="mt-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
				>
					{showUnchanged
						? `↑ Hide ${unchanged.length} unchanged`
						: `↓ Show ${unchanged.length} unchanged`}
				</button>
			)}
		</div>
	);
}
