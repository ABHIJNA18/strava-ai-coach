-- Adds an explicit completion timestamp for the first activity synchronization.

ALTER TABLE athletes
ADD COLUMN IF NOT EXISTS initial_sync_completed_at TIMESTAMPTZ NULL;