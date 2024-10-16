namespace ATCLib.Messages
{
    public class Roger(List<MessageToken> tokens) : MessagePayload(tokens)
    {
    }

    public class RogerMessagePayloadParser : MessagePayloadParser
    {
        public override MessagePayload? Parse(List<MessageToken> tokens)
        {
            if (tokens.Count == 0)
            {
                return new Roger(tokens);
            }

            var rogerIndex = MessagePayloadParser.FindPhraseIndex(tokens, [["roger"], ["affirmative"], ["copy"]]);
            if (rogerIndex == -1)
            {
                return null;
            }

            return new Roger(tokens);
        }
    }
}