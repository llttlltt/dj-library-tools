import { useAtom, useAtomValue } from "@effect-atom/atom-react";
import {
	Check,
	Loader2,
	RefreshCw,
	ShieldAlert,
	ShieldCheck,
} from "lucide-react";
import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { UpdateOverlay } from "@/components/ui/update-overlay";
import { runPromise } from "@/lib/runtime";
import { cn } from "@/lib/utils";
import { AppService } from "@/services";
import {
	checkForUpdates,
	fixPermissions,
	isCheckingUpdatesAtom,
	lastCheckedAtAtom,
	permissionStatusAtom,
	setUpdateInterval,
	showCheckSuccessAtom,
	updateConfigAtom,
	updateInfoAtom,
	versionAtom,
} from "@/store/system";

export default function SettingsView() {
	const [version] = useAtom(versionAtom);
	const [updateConfig] = useAtom(updateConfigAtom);
	const [permissionStatus] = useAtom(permissionStatusAtom);
	const [updateInfo] = useAtom(updateInfoAtom);
	const [isChecking] = useAtom(isCheckingUpdatesAtom);
	const [showSuccess] = useAtom(showCheckSuccessAtom);
	const lastCheckedAt = useAtomValue(lastCheckedAtAtom);

	const [updateStatus, setUpdateStatus] = useState<
		"idle" | "downloading" | "installing" | "complete" | "error"
	>("idle");
	const [error, setError] = useState<string>("");

	async function handleCheckUpdate() {
		await runPromise(checkForUpdates);
	}

	async function handleInstallUpdate() {
		if (!updateInfo) return;
		setUpdateStatus("downloading");
		try {
			const app = await runPromise(AppService);
			await runPromise(app.installUpdate());
			setUpdateStatus("complete");
		} catch (e) {
			setError(e instanceof Error ? e.message : String(e));
			setUpdateStatus("error");
		}
	}

	async function handleIntervalChange(value: string) {
		const hours = parseInt(value, 10);
		await runPromise(setUpdateInterval(hours));
	}

	async function handleFixPermissions() {
		await runPromise(fixPermissions);
	}

	const isDev = version.includes("-dev") || version === "v0.0.0";
	const isLatest = updateInfo !== null && !updateInfo.available;

	return (
		<div className="flex flex-col h-full">
			{/* Sticky Top Header Nav */}
			<div className="h-14 flex items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10">
				<span className="text-sm font-semibold">Settings</span>
				<div className="flex-1" />
			</div>

			{/* Scrollable Main Content Box */}
			<div className="flex-1 overflow-auto p-6">
				{updateStatus !== "idle" && (
					<UpdateOverlay
						status={updateStatus}
						version={updateInfo?.version}
						error={error}
					/>
				)}

				<div className="space-y-6">
					{/* Permissions Card */}
					<Card>
						<CardHeader className="border-b border-border pb-6">
							<CardTitle className="text-lg">Permissions</CardTitle>
							<CardDescription>
								Configure permissions required for the application's core
								background functions.
							</CardDescription>
						</CardHeader>
						<CardContent className="pt-6">
							<div className="flex items-center justify-between">
								<div className="space-y-1">
									<p className="text-sm font-medium">
										Installation Permissions
									</p>
									<p className="text-xs text-muted-foreground">
										Required for the application to apply future administrative
										updates automatically.
									</p>
								</div>
								<div className="flex items-center gap-3">
									{permissionStatus === "Healthy" ? (
										<Badge
											variant="outline"
											className="text-green-500 border-green-500/20 bg-green-500/5 py-1 px-2.5"
										>
											<ShieldCheck className="w-3.5 h-3.5 mr-1.5" /> Healthy
										</Badge>
									) : (
										<Badge
											variant="outline"
											className="text-yellow-500 border-yellow-500/20 bg-yellow-500/5 py-1 px-2.5"
										>
											<ShieldAlert className="w-3.5 h-3.5 mr-1.5" />{" "}
											{permissionStatus}
										</Badge>
									)}
									{permissionStatus !== "Healthy" && (
										<Button
											variant="outline"
											size="sm"
											onClick={handleFixPermissions}
										>
											Fix
										</Button>
									)}
								</div>
							</div>
						</CardContent>
					</Card>

					{/* Updates Card */}
					<Card>
						<CardHeader className="border-b border-border pb-6">
							<div className="flex items-center justify-between">
								<div className="space-y-1">
									<CardTitle className="text-lg">Application Updates</CardTitle>
									<CardDescription>
										Manage software updates and checking preferences.
									</CardDescription>
								</div>
								{updateInfo?.available && (
									<Badge
										variant="secondary"
										className="bg-green-500/10 text-green-500 border-green-500/20"
									>
										Update Available: {updateInfo.version}
									</Badge>
								)}
							</div>
						</CardHeader>
						<CardContent className="divide-y divide-border pt-6 space-y-6 [&>div:not(:first-child)]:pt-6">
							{/* Current Version & Update Check Row */}
							<div className="flex items-center justify-between">
								<div className="space-y-1">
									<p className="text-sm font-medium">Software Version</p>
									<div className="text-xs text-muted-foreground space-y-0.5">
										<div className="flex items-center gap-1.5">
											<span>
												Currently on{" "}
												<span className="font-semibold">{version}</span>.
											</span>
											{updateInfo?.available && (
												<span className="text-green-500">
													Newer version: {updateInfo.version}
												</span>
											)}
											{isLatest && !isDev && (
												<span className="text-emerald-500 font-medium inline-flex items-center gap-1">
													<Check className="w-3 h-3" /> You are running the
													latest version.
												</span>
											)}
											{isLatest && isDev && (
												<span className="text-muted-foreground italic">
													(Development version
													{updateInfo?.version &&
														updateInfo.version !== "v0.0.0" &&
														updateInfo.version !== "" &&
														` — Latest Release: ${updateInfo.version}`}
													)
												</span>
											)}
										</div>
										{lastCheckedAt && (
											<div className="flex items-center gap-1.5 opacity-60">
												<span>Last checked:</span>
												<span>{lastCheckedAt}</span>
											</div>
										)}
									</div>
								</div>
								<Button
									onClick={
										updateInfo?.available
											? handleInstallUpdate
											: handleCheckUpdate
									}
									disabled={isChecking || showSuccess}
									variant={
										updateInfo?.available
											? "default"
											: showSuccess
												? "outline"
												: "secondary"
									}
									size="sm"
									className={cn(
										"min-w-[140px]",
										showSuccess &&
											"text-emerald-500 border-emerald-500/20 bg-emerald-500/5 hover:bg-emerald-500/10",
									)}
								>
									{isChecking ? (
										<Loader2 className="w-4 h-4 mr-2 animate-spin" />
									) : showSuccess ? (
										<Check className="w-4 h-4 mr-2" />
									) : (
										<RefreshCw className="w-4 h-4 mr-2" />
									)}
									{updateInfo?.available
										? "Download Update"
										: showSuccess
											? "Latest"
											: "Check for Updates"}
								</Button>
							</div>

							{/* Frequency Row */}
							<div className="flex items-center justify-between">
								<div className="space-y-1">
									<p className="text-sm font-medium">Update Frequency</p>
									<p className="text-xs text-muted-foreground">
										Configure how often the app checks for releases in the
										background.
									</p>
								</div>
								<Select
									value={updateConfig?.check_interval_hour?.toString() || "168"}
									onValueChange={handleIntervalChange}
								>
									<SelectTrigger className="w-[140px]">
										<SelectValue placeholder="Select interval" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="24">Daily</SelectItem>
										<SelectItem value="168">Weekly</SelectItem>
									</SelectContent>
								</Select>
							</div>
						</CardContent>
					</Card>
				</div>
			</div>
		</div>
	);
}
