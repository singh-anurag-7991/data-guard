export interface ErrorDetail {
  rule_id: string;
  field: string;
  value: any;
  reason: string;
}

export interface ValidationResult {
  source_id: string;
  status: "PASS" | "FAIL";
  records_checked: number;
  rules_failed: number;
  errors?: ErrorDetail[];
  timestamp: string; // ISO string
}
