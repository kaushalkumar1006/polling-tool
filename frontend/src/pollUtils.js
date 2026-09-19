export function pollPercent(votes, total) {
  return total > 0 ? Math.round((votes / total) * 100) : 0;
}

export function canSubmitVote(optionId, submitting, voted) {
  return Boolean(optionId) && !submitting && !voted;
}
