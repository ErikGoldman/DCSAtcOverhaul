using System.Text.RegularExpressions;

namespace ATCLib.Messages
{
  public class RadioCheck(string messageString) : MessagePayload(messageString)
  {
  }

  public class RadioCheckMessagePayloadParser : MessagePayloadParser
  {
    private static readonly Regex RogerRegex = new Regex(@"(radio check|check radio|do you copy)", RegexOptions.IgnoreCase);

    public override MessagePayload? Parse(string messageString, List<MessageToken> tokens)
    {
      if (RogerRegex.IsMatch(messageString))
      {
        return new RadioCheck(messageString);
      }
      return null;
    }
  }
}