import { useAtom } from "@effect-atom/atom-react";
import { TreeFormatter } from "effect/ParseResult";
import { Loader2 } from "lucide-react";
import { useEffect, useState } from "react";
import { EndpointEditor } from "@/components/endpoint/EndpointEditor";
import { ResourceTable } from "@/components/resource/ResourceTable";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from "@/components/ui/sheet";
import { runPromise } from "@/lib/runtime";
import { AppService } from "@/services";
import { WailsDecodeError } from "@/services/errors";
import { connectionsAtom, loadConnections } from "@/store/connections";
import { loadProviders, providersAtom } from "@/store/providers";
import {
	filterWritableResources,
	findConnectionProvider,
} from "@/store/selection";
import type { QueryResult } from "@/types";

interface QueryTesterProps {
	open: boolean;
	onClose: () => void;
	initialConnectionID?: string;
	initialResource?: string;
	initialQuery?: string;
	isTarget?: boolean;
	onApply?: (query: string) => void;
}

const asQueryResult = (x: unknown) => x as QueryResult;

export function QueryTester({
	open,
	onClose,
	initialConnectionID,
	initialResource,
	initialQuery,
	isTarget,
	onApply,
}: QueryTesterProps) {
	const [connections] = useAtom(connectionsAtom);
	const [providers] = useAtom(providersAtom);

	const [connectionID, setConnectionID] = useState(initialConnectionID ?? "");
	const [resource, setResource] = useState(initialResource ?? "tracks");
	const [query, setQuery] = useState(initialQuery ?? "");
	const [result, setResult] = useState<QueryResult | null>(null);
	const [error, setError] = useState<unknown | null>(null);
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		runPromise(loadConnections);
		runPromise(loadProviders);
	}, []);

	// Automatically fix blank or invalid resource selections when connection or providers change
	useEffect(() => {
		if (connectionID && providers.length > 0) {
			const provider = findConnectionProvider(
				connectionID,
				connections,
				providers,
			);
			const availableResources = filterWritableResources(provider, !!isTarget);

			if (availableResources.length > 0) {
				const isValid = availableResources.some((r) => r.name === resource);
				if (!isValid) {
					setResource(availableResources[0].name);
				}
			}
		}
	}, [connectionID, providers, resource, isTarget, connections]);

	useEffect(() => {
		setConnectionID(initialConnectionID ?? "");
	}, [initialConnectionID]);

	useEffect(() => {
		setResource(initialResource ?? "tracks");
	}, [initialResource]);

	useEffect(() => {
		setQuery(initialQuery ?? "");
	}, [initialQuery]);

	// biome-ignore lint/correctness/useExhaustiveDependencies: intentionally watching value changes to clear stale results
	useEffect(() => {
		setResult(null);
		setError("");
	}, [connectionID, resource, query]);

	async function handleTest() {
		setLoading(true);
		setError(null);
		setResult(null);
		try {
			const app = await runPromise(AppService);
			const data = await runPromise(
				app.previewQuery(connectionID, resource, query),
			);
			setResult(asQueryResult(data));
		} catch (e) {
			setError(e);
		}
		setLoading(false);
	}

	const provider = findConnectionProvider(connectionID, connections, providers);
	const currentRes = provider?.resources.find((r) => r.name === resource);
	const isInvalidTarget = isTarget && currentRes && !currentRes.can_write;

	return (
		<Sheet open={open} onOpenChange={(o) => !o && onClose()}>
			<SheetContent className="flex flex-col h-full sm:max-w-xl md:max-w-2xl">
				<SheetHeader className="shrink-0 mb-4">
					<SheetTitle>Query Tester</SheetTitle>
					<SheetDescription>
						Test any query before using it in a Step.
					</SheetDescription>
				</SheetHeader>

				<div className="flex flex-col gap-6 flex-1 min-h-0">
					<div className="space-y-4 bg-secondary/20 p-4 rounded-xl border border-border/40 shrink-0">
						<EndpointEditor
							endpoint={{ connection_id: connectionID, resource, query }}
							connections={connections}
							providers={providers}
							isTarget={isTarget}
							onChange={(p) => {
								if (p.connection_id) setConnectionID(p.connection_id);
								if (p.resource) setResource(p.resource);
								if (p.query !== undefined) setQuery(p.query);
							}}
							layout="grid"
						/>

						{/* Actions Row */}
						<div className="flex gap-2 pt-1 border-t border-border/20">
							<Button
								type="button"
								size="sm"
								onClick={handleTest}
								disabled={loading || !connectionID}
								className="min-w-[80px]"
							>
								{loading ? (
									<>
										<Loader2 className="h-3.5 w-3.5 mr-1.5 animate-spin" />
										Testing…
									</>
								) : (
									"Test Query"
								)}
							</Button>

							{onApply && (
								<Button
									type="button"
									variant="secondary"
									size="sm"
									onClick={() => {
										onApply(query);
										onClose();
									}}
									disabled={!query || isInvalidTarget}
								>
									Use this query
								</Button>
							)}
						</div>
					</div>

					<div className="flex-1 min-h-0 flex flex-col">
						<QueryTesterResults result={result} error={error} />
					</div>
				</div>
			</SheetContent>
		</Sheet>
	);
}

// ── Shared sub-components used by both QueryTester (sheet) and QueryTesterView ──

interface ResultsProps {
	result: QueryResult | null;
	error: unknown | null;
}

export function QueryTesterResults({ result, error }: ResultsProps) {
	if (error) {
		let message = error instanceof Error ? error.message : String(error);

		if (error instanceof WailsDecodeError) {
			message = `${message}\n${TreeFormatter.formatError(error.parseError)}`;
		}

		return (
			<div className="p-4 rounded-xl border border-destructive/20 bg-destructive/5 text-sm text-destructive font-mono overflow-auto shrink-0 leading-relaxed whitespace-pre-wrap">
				<div className="font-semibold mb-1">Execution Error</div>
				{message}
			</div>
		);
	}
	if (result === null) return null;

	return (
		<div className="flex flex-col gap-3 flex-1 min-h-0">
			<div className="flex items-center justify-between">
				<span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
					Result Preview
				</span>
				<Badge
					variant="secondary"
					className="bg-emerald-500/10 text-emerald-500 border-emerald-500/20 py-0.5 px-2 text-xs font-medium"
				>
					Matched {result.count.toLocaleString()}{" "}
					{result.kind === "groups"
						? result.count !== 1
							? "items"
							: "item"
						: result.count !== 1
							? "tracks"
							: "track"}
				</Badge>
			</div>

			<ResourceTable result={result} />
		</div>
	);
}
