-- name: GetAvailableTranscriptionLanguages :many
SELECT
    tl.id,
    tl.code,
    tl.value,
    tl.country,
    tl.language,
    tl.boost_meeting_assistant,
    tl.boost_meeting_title,
    tl.boost_meeting_participants,
    tl.boost_meeting_participants_actual,
    tl.boost_template_custom_words,
    tl.boost_template_keyword_based_highlights
FROM
    transcription_languages tl
WHERE
    tl.available = 1
;
