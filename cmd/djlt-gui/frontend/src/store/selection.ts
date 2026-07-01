import type { Connection, ProviderInfo, ResourceInfo } from "@/types";

/**
 * selection.ts
 * Pure logic and helper functions for connection/resource selection.
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

export const findConnectionProvider = (
	connectionID: string,
	connections: Connection[],
	providers: ProviderInfo[],
): ProviderInfo | undefined => {
	const connection = connections.find((c) => c.id === connectionID);
	return providers.find((p) => p.name === connection?.provider);
};

export const getFirstValidResource = (
	connectionID: string,
	connections: Connection[],
	providers: ProviderInfo[],
	isTarget: boolean,
): string => {
	const provider = findConnectionProvider(connectionID, connections, providers);
	const valid = filterWritableResources(provider, isTarget);
	return valid[0]?.name ?? "tracks";
};

export const canConnectionSupportTarget = (
	connectionID: string,
	connections: Connection[],
	providers: ProviderInfo[],
): boolean => {
	const provider = findConnectionProvider(connectionID, connections, providers);
	return provider?.resources.some((r) => r.can_write) ?? true;
};
