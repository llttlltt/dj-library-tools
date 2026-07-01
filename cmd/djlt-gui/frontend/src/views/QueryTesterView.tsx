import { useAtom } from "@effect-atom/atom-react";
import { Loader2 } from "lucide-react";
import { useEffect, useState } from "react";
import { EndpointEditor } from "@/components/endpoint/EndpointEditor";
import { QueryTesterResults } from "@/components/query/QueryTester";
import { Button } from "@/components/ui/button";
import { runPromise } from "@/lib/runtime";
import { AppService } from "@/services";
import { connectionsAtom, loadConnections } from "@/store/connections";
import { loadProviders, providersAtom } from "@/store/providers";
import type { QueryResult } from "@/types";

const asQueryResult = (x: unknown) => x as QueryResult;

export default function QueryTesterView() {
	const [connections] = useAtom(connectionsAtom);
	const [providers] = useAtom(providersAtom);

	const [connectionID, setConnectionID] = useState("");
	const [resource, setResource] = useState("tracks");
	const [query, setQuery] = useState("");
	const [result, setResult] = useState<QueryResult | null>(null);
	const [error, setError] = useState<unknown | null>(null);
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		runPromise(loadConnections);
		runPromise(loadProviders);
	}, []);

	useEffect(() => {
		if (connections.length > 0 && !connectionID) {
			setConnectionID(connections[0].id);
		}
	}, [connections, connectionID]);

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

	return (
		<div className="flex flex-col h-full overflow-hidden">
			{/* Sticky Top Header Nav */}
			<div className="h-14 flex items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] shrink-0 z-10">
				<span className="text-sm font-semibold">Query Tester</span>
				<div className="flex-1" />
			</div>

			{/* Main Layout Container */}
			<div className="flex-1 p-6 flex flex-col min-h-0 bg-background">
				<div className="flex flex-col gap-6 h-full min-h-0">
					{/* Controls Box - Pins to Top */}
					<div className="shrink-0 space-y-4 bg-secondary/20 p-4 rounded-xl border border-border/40">
						<EndpointEditor
							endpoint={{ connection_id: connectionID, resource, query }}
							connections={connections}
							providers={providers}
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
						</div>
					</div>

					{/* Results / Error Panel - Only mounts if result exists or there is an active error */}
					{(result !== null || error !== null) && (
						<div className="flex-1 min-h-0 flex flex-col">
							<QueryTesterResults result={result} error={error} />
						</div>
					)}
				</div>
			</div>
		</div>
	);
}
