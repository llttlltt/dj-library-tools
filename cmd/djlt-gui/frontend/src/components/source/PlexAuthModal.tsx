import { Copy, ExternalLink, Loader2 } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { CheckPlexAuth, InitPlexAuth } from "../../../wailsjs/go/gui/App";

interface PlexAuthModalProps {
	open: boolean;
	onToken: (token: string) => void;
	onCancel: () => void;
}

export function PlexAuthModal({ open, onToken, onCancel }: PlexAuthModalProps) {
	const [authUrl, setAuthUrl] = useState("");
	const [pinId, setPinId] = useState(0);
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);
	const [copied, setCopied] = useState(false);
	const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);

	// Start PIN flow when opened
	useEffect(() => {
		if (!open) return;
		setError("");
		setAuthUrl("");
		setPinId(0);
		setLoading(true);

		(async () => {
			try {
				const challenge = await InitPlexAuth();
				setAuthUrl(challenge.url);
				setPinId(challenge.pin_id);
			} catch (e) {
				setError(String(e));
			}
			setLoading(false);
		})();
	}, [open]);

	// Poll for token once we have a PIN ID
	useEffect(() => {
		if (!pinId || !open) return;
		pollRef.current = setInterval(async () => {
			try {
				const token = await CheckPlexAuth(pinId);
				if (token) {
					clearInterval(pollRef.current ?? undefined);
					onToken(token);
				}
			} catch {
				// ignore transient poll errors
			}
		}, 2000);
		return () => clearInterval(pollRef.current ?? undefined);
	}, [pinId, open, onToken]);

	function copyUrl() {
		navigator.clipboard.writeText(authUrl).then(() => {
			setCopied(true);
			setTimeout(() => setCopied(false), 2000);
		});
	}

	function handleCancel() {
		clearInterval(pollRef.current ?? undefined);
		onCancel();
	}

	return (
		<Dialog
			open={open}
			onOpenChange={(o) => {
				if (!o) handleCancel();
			}}
		>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Authenticate with Plex</DialogTitle>
				</DialogHeader>

				<div className="flex flex-col gap-4 mt-2">
					{loading && (
						<div className="flex items-center gap-2 text-sm text-muted-foreground">
							<Loader2 className="h-4 w-4 animate-spin" />
							Requesting PIN…
						</div>
					)}

					{error && <p className="text-sm text-destructive">{error}</p>}

					{authUrl && (
						<>
							<p className="text-sm text-muted-foreground">
								Open the link below in your browser to authorise DJ Library
								Tools with your Plex account. This dialog will close
								automatically once authenticated.
							</p>

							<div className="flex items-center gap-2 p-3 rounded-md bg-muted/30 border border-border">
								<a
									href={authUrl}
									target="_blank"
									rel="noreferrer"
									className="text-xs text-blue-400 break-all flex-1 hover:underline"
								>
									{authUrl}
								</a>
								<ExternalLink className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
							</div>

							<div className="flex items-center gap-2 text-xs text-muted-foreground">
								<Loader2 className="h-3 w-3 animate-spin" />
								Waiting for authorisation…
							</div>

							<div className="flex justify-between gap-2 mt-1">
								<Button
									type="button"
									variant="outline"
									size="sm"
									onClick={copyUrl}
									className="gap-1.5"
								>
									<Copy className="h-3.5 w-3.5" />
									{copied ? "Copied!" : "Copy link"}
								</Button>
								<Button
									type="button"
									variant="ghost"
									size="sm"
									onClick={handleCancel}
								>
									Cancel
								</Button>
							</div>
						</>
					)}
				</div>
			</DialogContent>
		</Dialog>
	);
}
