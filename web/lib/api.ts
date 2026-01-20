import { ValidationResult } from "./types";

const API_BASE_URL = "http://localhost:8080";

export async function getRecentRuns(sourceID?: string): Promise<ValidationResult[]> {
    const url = new URL(`${API_BASE_URL}/api/runs`);
    if (sourceID) {
        url.searchParams.set("source_id", sourceID);
    }

    try {
        const res = await fetch(url.toString(), { cache: "no-store" });
        if (!res.ok) {
            throw new Error(`Failed to fetch runs: ${res.statusText}`);
        }
        const data = await res.json();
        return data || [];
    } catch (error) {
        console.error("API Fetch Error:", error);
        return [];
    }
}
