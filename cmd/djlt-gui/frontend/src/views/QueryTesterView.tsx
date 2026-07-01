import { useAtom } from "@effect-atom/atom-react";
import { useEffect, useState } from "react";
import {
	QueryTesterControls,
	QueryTesterResults,
} from "@/components/query/QueryTester";
import { runtime } from "@/lib/runtime";
import { loadProviders, providersAtom } from "@/store/providers";
import { loadSources, sourcesAtom } from "@/store/sources";
import type { QueryResult } from "@/types";
import { PreviewQuery } from "../../wailsjs/go/gui/App";

const asQueryResult = (x: unknown) => x as QueryResult;

export default function QueryTesterView() {
	const [sources] = useAtom(sourcesAtom);
	const [providers] = useAtom(providersAtom);

	const [sourceID, setSourceID] = useState("");
	const [resource, setResource] = useState("tracks");
	const [query, setQuery] = useState("");
	const [result, setResult] = useState<QueryResult | null>(null);
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		runtime.runPromise(loadSources);
		runtime.runPromise(loadProviders);
	}, []);

	useEffect(() => {
		if (sources.length > 0 && !sourceID) {
			setSourceID(sources[0].id);
		}
	}, [sources, sourceID]);

	// biome-ignore lint/correctness/useExhaustiveDependencies: intentionally watching value changes to clear stale results
	useEffect(() => {
		setResult(null);
		setError("");
	}, [sourceID, resource, query]);

	async function handleTest() {
		setLoading(true);
		setError("");
		setResult(null);
		try {
			setResult(asQueryResult(await PreviewQuery(sourceID, resource, query)));
		} catch (e) {
			setError(String(e));
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
					<div className="shrink-0">
						<QueryTesterControls
							sources={sources}
							providers={providers}
							sourceID={sourceID}
							resource={resource}
							query={query}
							loading={loading}
							onSourceID={setSourceID}
							onResource={setResource}
							onQuery={setQuery}
							onTest={handleTest}
						/>
					</div>

					{/* Results / Error Panel - Only mounts if result exists or there is an active error */}
					{(result !== null || error) && (
						<div className="flex-1 min-h-0 flex flex-col">
							<QueryTesterResults result={result} error={error} />
						</div>
					)}
				</div>
			</div>
		</div>
	);
}
