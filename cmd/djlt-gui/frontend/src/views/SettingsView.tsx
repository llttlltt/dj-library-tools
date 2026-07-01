import { Loader2, RefreshCw, ShieldAlert, ShieldCheck } from "lucide-react";
import { useEffect, useState } from "react";
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
import {
	CheckForUpdate,
	FixPermissions,
	GetPermissionStatus,
	GetUpdateConfig,
	GetVersion,
	InstallUpdate,
	SetUpdateInterval,
} from "../../wailsjs/go/gui/App";
import type { config, update } from "../../wailsjs/go/models";

export default function SettingsView() {
	const [version, setVersion] = useState("v0.0.0");
	const [updateConfig, setUpdateConfig] = useState<config.UpdateConfig | null>(
		null,
	);
	const [permissionStatus, setPermissionStatus] =
		useState<string>("Checking...");
	const [isChecking, setIsChecking] = useState(false);
	const [updateInfo, setUpdateInfo] = useState<update.UpdateInfo | null>(null);
	const [updateStatus, setUpdateStatus] = useState<
		"idle" | "downloading" | "installing" | "complete" | "error"
	>("idle");
	const [error, setError] = useState<string>("");

	useEffect(() => {
		GetVersion().then(setVersion);
		GetUpdateConfig().then(setUpdateConfig);
		GetPermissionStatus().then(setPermissionStatus);
	}, []);

	async function handleCheckUpdate() {
		setIsChecking(true);
		try {
			const info = await CheckForUpdate(true);
			setUpdateInfo(info);
		} finally {
			setIsChecking(false);
		}
	}

	async function handleInstallUpdate() {
		if (!updateInfo) return;
		setUpdateStatus("downloading");
		try {
			await InstallUpdate();
			setUpdateStatus("complete");
		} catch (e) {
			setError(e instanceof Error ? e.message : String(e));
			setUpdateStatus("error");
		}
	}

	async function handleIntervalChange(value: string) {
		const hours = parseInt(value, 10);
		await SetUpdateInterval(hours);
		setUpdateConfig((prev) =>
			prev ? { ...prev, check_interval_hour: hours } : null,
		);
	}

	async function handleFixPermissions() {
		try {
			await FixPermissions();
			const status = await GetPermissionStatus();
			setPermissionStatus(status);
		} catch (e) {
			console.error("Failed to fix permissions:", e);
		}
	}

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
									<p className="text-xs text-muted-foreground">
										Currently on{" "}
										<span className="font-semibold">{version}</span>.{" "}
										{updateInfo?.available
											? `A newer version (${updateInfo.version}) is ready to download.`
											: "You are running the latest version."}
									</p>
								</div>
								<Button
									onClick={
										updateInfo?.available
											? handleInstallUpdate
											: handleCheckUpdate
									}
									disabled={isChecking}
									variant={updateInfo?.available ? "default" : "secondary"}
									size="sm"
								>
									{isChecking ? (
										<Loader2 className="w-4 h-4 mr-2 animate-spin" />
									) : (
										<RefreshCw className="w-4 h-4 mr-2" />
									)}
									{updateInfo?.available
										? "Download Update"
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
