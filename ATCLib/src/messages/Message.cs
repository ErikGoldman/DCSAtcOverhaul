namespace ATCLib
{
  public abstract class MessagePayloadParser
  {
    public abstract MessagePayload? Parse(string messageString, List<MessageToken> tokens);
    public static int FindPhraseIndex(List<MessageToken> tokens, string[][] phrases)
    {
      for (int i = 0; i < tokens.Count; i++)
      {
        foreach (var phrase in phrases)
        {
          if (phrase.Length > tokens.Count - i)
          {
            continue;
          }

          if (tokens[i].Content.Equals(phrase[0], StringComparison.CurrentCultureIgnoreCase))
          {
            bool match = true;
            for (int j = 1; j < phrase.Length; j++)
            {
              if (!tokens[i + j].Content.Equals(phrase[j], StringComparison.CurrentCultureIgnoreCase))
              {
                match = false;
                break;
              }
            }

            if (match)
            {
              return i;
            }
          }
        }
      }

      return -1;
    }
  }

  public abstract class MessagePayload
  {
    public string RawContent { get; }

    public MessagePayload(string rawContent)
    {
      RawContent = rawContent;
    }
  }

  public class EmptyMessagePayload(string rawContent) : MessagePayload(rawContent)
  {
  }

  public class EmptyMessagePayloadParser : MessagePayloadParser
  {
    public override MessagePayload? Parse(string messageString, List<MessageToken> tokens)
    {
      return new EmptyMessagePayload(messageString);
    }
  }

  public class MessageParseException : Exception
  {
    public MessageParseException(string message) : base(message) { }
  }

  public readonly struct MessageToken(bool isCallsign, string content)
  {
    public bool IsCallsign { get; } = isCallsign;
    public string Content { get; } = content;
  }

  public class Message
  {
    public string RawMessage { get; }

    public ICommunicator? From { get; }
    public ICommunicator? To { get; }
    public MessagePayload Payload { get; }

    public Message(string rawMessage, ICommunicator? from, ICommunicator? to, MessagePayload payload)
    {
      RawMessage = rawMessage;
      From = from;
      To = to;
      Payload = payload;
    }

    // this function tokenizes the message by extracting the callsigns and the message content separately
    public static List<MessageToken> SplitMessage(string rawMessage, ActiveCommunicatorList communicators)
    {
      if (rawMessage == "")
      {
        return [];
      }

      string[] inTokens = rawMessage.Split(' ');
      if (inTokens.Length == 0)
      {
        return [];
      }

      List<MessageToken> outTokens = [];
      int currentCallsignParsePositionStart = 0;
      int currentCallsignParsePositionEnd = 0;

      while (currentCallsignParsePositionStart < inTokens.Length)
      {
        var callsignSoFar = string.Join(" ", inTokens[currentCallsignParsePositionStart..(currentCallsignParsePositionEnd + 1)]);

        if (communicators.GetByName(callsignSoFar) != null)
        {
          outTokens.Add(new MessageToken(true, callsignSoFar));
          currentCallsignParsePositionStart = currentCallsignParsePositionEnd + 1;
          currentCallsignParsePositionEnd = currentCallsignParsePositionStart;
        }
        else if (communicators.PrefixExists(callsignSoFar))
        {
          currentCallsignParsePositionEnd++;
        }
        else
        {
          outTokens.Add(new MessageToken(false, inTokens[currentCallsignParsePositionStart]));
          currentCallsignParsePositionStart++;
          currentCallsignParsePositionEnd = currentCallsignParsePositionStart;
        }
      }

      return outTokens;
    }

    public static Message? Parse(string rawMessage, ATCState state)
    {
      List<MessageToken> tokens = SplitMessage(rawMessage, state.ActiveCommunicators);
      ICommunicator? maybeFrom = null;
      ICommunicator? maybeTo = null;
      int contentStart = 0, contentEnd = tokens.Count;

      // format possibilities:
      // <callsign> <callsign> <message>
      // <callsign> <message>
      // <message> <callsign>
      // <message>


      if (tokens.Count == 0)
      {
        Log.WriteLine("Message.Parse: Empty message");
        return null;
      }
      else if (tokens.Count == 1)
      {
        if (tokens[0].IsCallsign)
        {
          maybeFrom = state.ActiveCommunicators.GetByName(tokens[0].Content);
          contentStart = 1;
        }
        else
        {
          contentStart = 0;
        }
      }
      else
      {
        if (tokens[0].IsCallsign && !tokens[1].IsCallsign)
        {
          // check if it's to at the beginning, from at the end
          if (tokens[^1].IsCallsign)
          {
            maybeTo = state.ActiveCommunicators.GetByName(tokens[0].Content);
            maybeFrom = state.ActiveCommunicators.GetByName(tokens[^1].Content);
            contentStart = 1;
            contentEnd = tokens.Count - 1;
          }
          else
          {
            maybeFrom = state.ActiveCommunicators.GetByName(tokens[0].Content);
            contentStart = 1;
          }
        }
        else if (tokens[0].IsCallsign && tokens[1].IsCallsign)
        {
          maybeTo = state.ActiveCommunicators.GetByName(tokens[0].Content);
          maybeFrom = state.ActiveCommunicators.GetByName(tokens[1].Content);
          contentStart = 2;
        }

        if (maybeFrom == null && tokens[^1].IsCallsign)
        {
          maybeFrom = state.ActiveCommunicators.GetByName(tokens[^1].Content);
          contentEnd = tokens.Count - 1;
        }
      }

      var messageString = string.Join(" ", tokens[contentStart..contentEnd].Select(t => t.Content));
      var messageTokens = tokens[contentStart..contentEnd];

      foreach (var payloadParser in state.MessagePayloadTypes)
      {
        var payload = payloadParser.Parse(messageString, messageTokens);
        if (payload != null)
        {
          return new Message(rawMessage, maybeFrom, maybeTo, payload);
        }
      }

      if (maybeFrom != null)
      {
        return new Message(rawMessage, maybeFrom, maybeTo, new EmptyMessagePayload(messageString));
      }

      return null;
    }
  }
}