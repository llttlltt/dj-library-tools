import type { ProviderInfo, ResourceInfo, Source } from "@/types";

/**
 * selection.ts
 * Pure logic and helper functions for source/resource selection.
 * These will be consumed by both Atoms and Components.
 */

export const filterWritableResources = (
	provider: ProviderInfo | undefined,
	isTarget: boolean,
): ResourceInfo[] => {
	if (!provider) return [];
	return provider.resources.filter((r) => {
		if (!isTarget) return true;
		return r.can_write;
	});
};

export const findSourceProvider = (
	sourceID: string,
	sources: Source[],
	providers: ProviderInfo[],
): ProviderInfo | undefined => {
	const source = sources.find((s) => s.id === sourceID);
	return providers.find((p) => p.name === source?.provider);
};

export const getFirstValidResource = (
	sourceID: string,
	sources: Source[],
	providers: ProviderInfo[],
	isTarget: boolean,
): string => {
	const provider = findSourceProvider(sourceID, sources, providers);
	const valid = filterWritableResources(provider, isTarget);
	return valid[0]?.name ?? "tracks";
};

export const canSourceSupportTarget = (
	sourceID: string,
	sources: Source[],
	providers: ProviderInfo[],
): boolean => {
	const provider = findSourceProvider(sourceID, sources, providers);
	return provider?.resources.some((r) => r.can_write) ?? true;
};
