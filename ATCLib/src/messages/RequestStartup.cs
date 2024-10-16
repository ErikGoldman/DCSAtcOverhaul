
using System.Text.RegularExpressions;

namespace ATCLib.Messages
{
  public class RequestStartup(string messageString, List<MessageToken> otherPlanes) : MessagePayload(messageString)
  {
    public List<MessageToken> OtherPlanes { get; } = otherPlanes;
  }

  public class RequestStartupMessagePayloadParser : MessagePayloadParser
  {
    private static readonly Regex RequestStartupRegex = new Regex(@"request(ing)? (engine )?(start|power)( )?(up)?( (and|with) weather( info(rmation)?)?)?", RegexOptions.IgnoreCase);

    public override MessagePayload? Parse(string messageString, List<MessageToken> tokens)
    {
      if (RequestStartupRegex.IsMatch(messageString))
      {
        return new RequestStartup(messageString, tokens.Where(t => t.IsCallsign).ToList());
      }
      return null;
    }
  }
}