using System.Text.RegularExpressions;

namespace ATCLib.Messages
{
  public class SayAgain(string messageString) : MessagePayload(messageString)
  {
  }

  public class SayAgainMessagePayloadParser : MessagePayloadParser
  {
    private static readonly Regex SayAgainRegex = new Regex(@"(say ((it|that) )?again|repeat)", RegexOptions.IgnoreCase);

    public override MessagePayload? Parse(string messageString, List<MessageToken> tokens)
    {
      if (SayAgainRegex.IsMatch(messageString))
      {
        return new SayAgain(messageString);
      }

      return null;
    }
  }
}