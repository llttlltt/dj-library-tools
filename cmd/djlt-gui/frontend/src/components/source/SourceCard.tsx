import { ChevronRight, Trash2 } from "lucide-react";
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
			<CardHeader className="flex-row items-center justify-between py-3 px-4 gap-0 space-y-0">
				<div className="flex items-center gap-3 min-w-0 flex-1 mr-2">
					<CardTitle className="text-sm font-medium truncate">
						{s.name}
					</CardTitle>
					<Badge
						variant={
							s.provider === "rb"
								? "sync"
								: s.provider === "plex"
									? "fix"
									: "edit"
						}
						className="shrink-0"
					>
						{PROVIDER_LABELS[s.provider] ?? s.provider}
					</Badge>
				</div>
				<div className="flex items-center gap-1">
					<Button
						type="button"
						variant="ghost"
						size="icon"
						className="h-8 w-8 shrink-0 hover:bg-secondary"
						onClick={(e) => {
							e.stopPropagation();
							onDelete(s);
						}}
					>
						<Trash2 className="h-4 w-4 text-muted-foreground hover:text-destructive" />
					</Button>
					<div className="h-4 w-px bg-border mx-1" />
					<ChevronRight className="h-4 w-4 text-muted-foreground/60" />
				</div>
			</CardHeader>
			{subtitle && (
				<CardContent className="py-0 pb-3 px-4 pr-16 -mt-1">
					<p className="text-xs text-muted-foreground truncate">{subtitle}</p>
				</CardContent>
			)}
		</Card>
	);
}
