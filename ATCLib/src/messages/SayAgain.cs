namespace ATCLib.Messages
{
  public class SayAgain(List<MessageToken> tokens) : MessagePayload(tokens)
  {
  }

  public class SayAgainMessagePayloadParser : MessagePayloadParser
  {
    public override MessagePayload? Parse(List<MessageToken> tokens)
    {
      var sayAgainIndex = MessagePayloadParser.FindPhraseIndex(tokens, [
        ["say", "again"], ["repeat"], ["say", "that", "again"]
      ]);
      if (sayAgainIndex == -1)
      {
        return null;
      }

      return new SayAgain(tokens);
    }
  }
}