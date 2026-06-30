import { Pencil, Trash2 } from "lucide-react";
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
	return (
		<Card>
			<CardHeader className="flex-row items-center justify-between py-3 px-4 gap-0">
				<div className="flex items-center gap-3">
					<CardTitle className="text-sm">{s.name}</CardTitle>
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
				<div className="flex items-center gap-0.5">
					<Button
						type="button"
						variant="ghost"
						size="icon"
						onClick={() => onEdit(s)}
					>
						<Pencil className="h-3.5 w-3.5 text-muted-foreground" />
					</Button>
					<Button
						type="button"
						variant="ghost"
						size="icon"
						onClick={() => onDelete(s)}
					>
						<Trash2 className="h-3.5 w-3.5 text-muted-foreground" />
					</Button>
				</div>
			</CardHeader>
			<CardContent className="py-0 pb-3 px-4">
				<p className="text-xs text-muted-foreground font-mono">{s.id}</p>
				{s.config?.file_path && (
					<p className="text-xs text-muted-foreground mt-0.5 truncate">
						{s.config.file_path}
					</p>
				)}
			</CardContent>
		</Card>
	);
}
