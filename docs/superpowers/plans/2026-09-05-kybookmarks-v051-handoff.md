**Repo:** kybookmarks-server
**PR:** #20 — https://github.com/Busness-app/kybookmarks-server/pull/20
**Worktree:** /home/yoshi/busness.app/kybookmarks-server/.claude/worktrees/drill-v051 (branch fix/recovery-drill-v051, clean)

# Post 291 execution hand-off — 2026-09-05

Owner: Usagi / GPT-6 / basalt.

## Done

Implemented and pushed the v0.5.1 opened-manifest migration, fixed product-required file/schema/admin validation, confined member access and read-only SQLite verification. HTTP and CLI now share a per-directory advisory lock before collection and scratch sweeping. Existing token settings, labels, ciphertext and key pins remain compatible. Added the synthetic v0.5.0 fixture, negative validation tests, subprocess locking tests, old-pairing deposit/unpair coverage and collected-instance restore coverage. Updated root AGENTS.md, README and RESTORE.md. The older master checkout, pre-existing untracked plans and locked deposit worktree were preserved.

Exact PR head: `850c6ce68a7a86377ca1126ca2a40f13aa77de61`.

- GitHub backend and frontend CI passed on that head.
- Autonomous security reviewer posted `verdict=cleared` on that same head: https://github.com/Busness-app/kybookmarks-server/pull/20#issuecomment-5552475888 . No blocking findings.
- Local Go build/vet/full race suite, frontend build and Docker build passed. All 42 ablations caught their intended failures. The final test-only commit also passed focused race tests and the full GitHub gates.
- A temporary forbidden `capsule.Open` call made the decrypt guard fail; removed it and the guard passed.
- The v0.5.0 fixture opens under v0.5.1 and drives a deposit without re-pairing. Wrong key/label fails. Unpair preserves receipts, local copies and pin.
- Synthetic collected-instance restore preserves all collected bytes and 0600 modes, the WAL canary, pairing readability and audit startup. Wrong service, insufficient shares, another ceremony's shares, damaged capsule and nonempty target are refused.

## Disposable Docker/TLS proof

Used an isolated KyBookmarks Docker container and a private TLS protocol fixture, not the production KyRecovery store. Exercised setup, pin/refused second key, local-only backup, HTTP and CLI drills including a held-lock conflict, schedule bounds/off, real HTTPS claim/deposit with explicit service name, refused mismatching receipt digest, identical local/remote bytes, restart without re-pairing, and Docker restore with synthetic shares on stdin. Booted the restored instance using separate data/config mounts and deposited with its restored pairing. Only one pairing claim occurred across all three deposits. Eight files restored; 40 drill checks passed. Unpair retained the key pin, receipt and local copies.

The proof image was built at `2632de586fe81f543ed1b0f4a74cdbccd798017a`. The only subsequent change in `850c6ce` is stronger tests in `internal/backup/drill_test.go`; production sources are identical. Image and nonsecret receipts:

```json
{
  "image_id": "sha256:73bed9a95ac9d60503ae3bce56dc494f104e8774f6975754079421643254081d",
  "fixture": "disposable Docker KyBookmarks + private TLS protocol fixture; not production KyRecovery",
  "key_id": "bd085c4c1bbd334c8fc498436578f5ed11c165c3cc50c8b721ae3e797f30f2ec",
  "claims": 1,
  "deposits": [
    {
      "capsule_id": "cap-KyBookmarks-1788617820131698503",
      "digest": "706c4eb55fabf095414715c0589a497f28c3d6dacab48fee499f1725911ec4d9",
      "size_bytes": 13195,
      "deposited_at": "2026-09-05T14:17:00.134893+00:00"
    },
    {
      "capsule_id": "cap-KyBookmarks-1788617820290777720",
      "digest": "3547cd2dfe1a5b86b510c7b7165609f19a775efe81ebfa3eb3f801c8a8024a2e",
      "size_bytes": 13846,
      "deposited_at": "2026-09-05T14:17:00.300336+00:00"
    },
    {
      "capsule_id": "cap-KyBookmarks-1788617822148277718",
      "digest": "66d596f2504a49d31f1ed50b76296c65861baec136b50d74e4af49150d75e453",
      "size_bytes": 14039,
      "deposited_at": "2026-09-05T14:17:02.156994+00:00"
    }
  ],
  "drill_checks": 40,
  "restore_files": 8,
  "checks": [
    "HTTP pin, refused second key",
    "local-only capsule 0600",
    "HTTP/CLI drill and busy conflict",
    "schedule bounds/off",
    "TLS claim and matching service_name",
    "digest mismatch refused",
    "same bytes local/remote",
    "restart deposit without re-pairing",
    "Docker restore stdin; wrong service, insufficient shares and nonempty target refused",
    "restored container booted and deposited without re-pairing",
    "unpair retained pin/receipt/local copies"
  ]
}
```

All disposable containers, extracted data, TLS private keys and synthetic share files were removed after the proof. The reviewer was invoked directly for this repository because the local pr-reviewer user service was not installed; its configuration was not changed. Its attempt to post to the HTTP board URL received 301, so it saved the detailed report locally at `/home/yoshi/.local/state/pr-reviewer/Busness-app-kybookmarks-server-20-850c6ce68a7a86377ca1126ca2a40f13aa77de61.md`. The cleared verdict is on GitHub.

## Left / exact input needed

PR #20 is ready for human merge; it has not been merged. Actual deployment proof remains blocked on the KyBookmarks deployment hostname/SSH alias and deployment path/access. This machine has KyRecovery running but no KyBookmarks container, and no remote Docker context is configured. The user was asked for the host during execution; no answer had arrived at hand-off time.

With that location available, inspect the existing deployment, preserve its volumes/issuer/keys/pairing, deploy the agreed reviewed revision, verify readiness/DNS/same pinned key and deposit without re-pairing. Record deployed SHA/image, capsule ID, digest, receipt time and local-copy result. Use `KYBOOKMARKS_BACKUP_ALLOW_PRIVATE_RECOVERY` and `KYBOOKMARKS_DNS` only for the intended homelab; source compose deployments require `up -d --build` with the LAN DNS override. Do not reset a pairing to make the proof pass. If the instance was never paired, label that first-pair proof instead.

## Careful

The disposable proof is not a production deposit or a real custodian-card restore. Real shares stay on the custodians' local stdin, never chat/argv/the board. A crash can leave cleartext drill residue until the library's next stale sweep; normal-return cleanup was verified. Keep the permanent lock inode and the existing decrypt-guard exemption. Code and CI/review are complete; deployment proof is a separate remaining stage. Check the board before resuming and preserve newer claims.
