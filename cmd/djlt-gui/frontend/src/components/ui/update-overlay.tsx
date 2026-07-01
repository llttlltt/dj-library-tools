import { Loader2, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";

interface UpdateOverlayProps {
	status: "downloading" | "installing" | "complete" | "error";
	version?: string;
	error?: string;
}

export function UpdateOverlay({ status, version, error }: UpdateOverlayProps) {
	return (
		<div className="fixed inset-0 z-[100] flex items-center justify-center bg-background/80 backdrop-blur-sm animate-in fade-in duration-300">
			<div className="w-full max-w-md p-6 bg-card border border-border rounded-xl shadow-2xl space-y-6 text-center">
				<div className="flex justify-center">
					<div className="p-3 bg-accent rounded-full">
						<RefreshCw
							className={`w-8 h-8 ${status !== "complete" && status !== "error" ? "animate-spin" : ""}`}
						/>
					</div>
				</div>

				<div className="space-y-2">
					<h2 className="text-2xl font-bold tracking-tight">
						{status === "downloading" && "Downloading Update"}
						{status === "installing" && "Installing Update"}
						{status === "complete" && "Update Complete"}
						{status === "error" && "Update Failed"}
					</h2>
					<p className="text-muted-foreground">
						{status === "downloading" &&
							`Fetching DJ Library Tools ${version}...`}
						{status === "installing" &&
							"Applying changes to the application..."}
						{status === "complete" &&
							"The application has been updated. Please restart to apply changes."}
						{status === "error" &&
							(error || "An unexpected error occurred during the update.")}
					</p>
				</div>

				{status === "complete" && (
					<Button
						className="w-full h-11 text-base font-medium"
						onClick={() => window.location.reload()}
					>
						Restart Now
					</Button>
				)}

				{status === "error" && (
					<Button
						className="w-full h-11 text-base font-medium"
						onClick={() => window.location.reload()}
					>
						Dismiss
					</Button>
				)}

				{(status === "downloading" || status === "installing") && (
					<div className="flex items-center justify-center gap-2 text-sm text-muted-foreground">
						<Loader2 className="w-4 h-4 animate-spin" />
						Please do not close the app
					</div>
				)}
			</div>
		</div>
	);
}
