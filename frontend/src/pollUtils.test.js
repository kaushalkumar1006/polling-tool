import { describe, expect, it } from 'vitest';
import { canSubmitVote, pollPercent } from './pollUtils';

describe('poll result helpers', () => {
  it('calculates rounded percentages and avoids zero division', () => {
    expect(pollPercent(2, 3)).toBe(67);
    expect(pollPercent(0, 0)).toBe(0);
  });

  it('only enables voting with a selected option and idle state', () => {
    expect(canSubmitVote('', false, false)).toBe(false);
    expect(canSubmitVote('option', true, false)).toBe(false);
    expect(canSubmitVote('option', false, true)).toBe(false);
    expect(canSubmitVote('option', false, false)).toBe(true);
  });
});
