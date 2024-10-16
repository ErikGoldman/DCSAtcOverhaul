using System.Text.RegularExpressions;

namespace ATCLib.Messages
{
    public class Roger(string messageString) : MessagePayload(messageString)
    {
    }

    public class RogerMessagePayloadParser : MessagePayloadParser
    {
        private static readonly Regex RogerRegex = new Regex(@"(roger|affirmative|copy)", RegexOptions.IgnoreCase);

        public override MessagePayload? Parse(string messageString, List<MessageToken> tokens)
        {
            if (tokens.Count == 0)
            {
                return new Roger(messageString);
            }

            if (RogerRegex.IsMatch(messageString))
            {
                return new Roger(messageString);
            }
            return null;
        }
    }
}