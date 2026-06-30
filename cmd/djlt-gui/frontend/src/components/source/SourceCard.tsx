import { Trash2 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Source } from "@/types";

const PROVIDER_LABELS: Record<string, string> = {
	rb: "Rekordbox",
	m3u: "M3U",
	plex: "Plex",
};

interface SourceCardProps {
	source: Source;
	onEdit: (s: Source) => void;
	onDelete: (s: Source) => void;
}

export function SourceCard({ source: s, onEdit, onDelete }: SourceCardProps) {
	const subtitle =
		s.config?.file_path || s.config?.host
			? (s.config.file_path ?? s.config.host)
			: null;

	return (
		<Card
			className="cursor-pointer hover:border-border/80 transition-colors"
			onClick={() => onEdit(s)}
		>
			<CardHeader className="flex-row items-center justify-between py-3 px-4 gap-0">
				<div className="flex items-center gap-3 min-w-0">
					<CardTitle className="text-sm truncate">{s.name}</CardTitle>
					<Badge
						variant={
							s.provider === "rb"
								? "sync"
								: s.provider === "plex"
									? "fix"
									: "edit"
						}
					>
						{PROVIDER_LABELS[s.provider] ?? s.provider}
					</Badge>
				</div>
				<Button
					type="button"
					variant="ghost"
					size="icon"
					className="shrink-0 ml-2"
					onClick={(e) => {
						e.stopPropagation();
						onDelete(s);
					}}
				>
					<Trash2 className="h-3.5 w-3.5 text-muted-foreground" />
				</Button>
			</CardHeader>
			{subtitle && (
				<CardContent className="py-0 pb-3 px-4">
					<p className="text-xs text-muted-foreground truncate">{subtitle}</p>
				</CardContent>
			)}
		</Card>
	);
}
